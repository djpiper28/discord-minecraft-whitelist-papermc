package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Goscord/goscord/discord"
	"github.com/Goscord/goscord/discord/embed"
	"github.com/Goscord/goscord/gateway"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

// Bot database model - a "singleton" for the server
type GuildSettings struct {
	gorm.Model
	// Guild Id
	ID string `gorm:"primaryKey"`
	// This is the admin role (@Committe 2022/23)
	AdminRole string
	// Role to get access to the system (@Student)
	AccessRole string
	// A master toggle to turn off user registrations
	AllowUserRegistrations bool
	MaxAccountsPerUser     int64
}

type DiscordUser struct {
	// This is a cache of whether or not the user has the admin role
	HasAdminRole bool `gorm:"index"`
	// This is whether the user and, all their accounts are banned
	Banned                bool                   `gorm:"index"`
	DiscordUserID         string                 `gorm:"primaryKey"`
	DiscordMinecraftUsers []DiscordMinecraftUser `gorm:"foreignKey:DiscordUserID"`
}

type DiscordMinecraftUser struct {
	DiscordUserID   string    `gorm:"foriegnKey:DiscordUser.DiscordUserID,index,unique,composite:discord_minecraft_user"`
	MinecraftUserID uuid.UUID `gorm:"foreignKey:MinecraftUser.Id,index,unique,composite:discord_minecraft_user"`
	Verified        bool
}

type MinecraftUser struct {
	Id                 uuid.UUID `gorm:"primaryKey"`
	Username           string
	LastLoginTime      time.Time
	LastX              float32
	LastY              float32
	LastZ              float32
	LastIpAddress      pgtype.Inet `gorm:"type:inet"`
	LastChunkImage     []byte
	LastSkinImage      []byte
	VerificationNumber int64
}

func reportMigrateError(err error) {
	if err != nil {
		log.Print(err)
	}
}

func AutoMigrateModel() {
	reportMigrateError(db.AutoMigrate(&MinecraftUser{}))
	reportMigrateError(db.AutoMigrate(&DiscordMinecraftUser{}))
	reportMigrateError(db.AutoMigrate(&DiscordUser{}))
	reportMigrateError(db.AutoMigrate(&GuildSettings{}))
}

// Helper function to set IP addresses, probably won't be used lmao
func SetInet(ip string) pgtype.Inet {
	var inet pgtype.Inet
	inet.Set(ip)
	return inet
}

type Context struct {
	client      *gateway.Session
	interaction *discord.Interaction
}

type Command interface {
	Name() string
	Description() string
	Category() string
	Options() []*discord.ApplicationCommandOption
	Execute(ctx *Context) bool
}

func Register(cmd Command, client *gateway.Session, commands map[string]Command) {
	appCmd := &discord.ApplicationCommand{
		Name:        cmd.Name(),
		Type:        discord.ApplicationCommandChat,
		Description: cmd.Description(),
		Options:     cmd.Options(),
	}

	_, err := client.Application.RegisterCommand(client.Me().Id, "", appCmd)
	if err != nil {
		log.Printf("Error registering command '%s' - %s", cmd.Name(), err)
	} else {
		log.Printf("Registered command '%s'", cmd.Name())
	}
	commands[cmd.Name()] = cmd
}

func ThemeEmbed(e *embed.Builder, ctx *Context) {
	e.SetFooter(ctx.client.Me().Username, ctx.client.Me().AvatarURL())
	e.SetColor(embed.Green)
	e.SetThumbnail(ctx.interaction.Member.User.AvatarURL())
}

func SendError(message string, ctx *Context) {
	e := embed.NewEmbedBuilder()

	e.SetTitle("An Error Occurred During Your Command")
	e.SetDescription(message)
	ThemeEmbed(e, ctx)

	// Send response
	ctx.client.Interaction.CreateResponse(ctx.interaction.Id,
		ctx.interaction.Token,
		&discord.InteractionCallbackMessage{Embeds: []*embed.Embed{e.Embed()},
			Flags: discord.MessageFlagEphemeral})
}

func SendAdminPermissionsError(gs GuildSettings, ctx *Context) {
	SendError(fmt.Sprintf("You require the <@&%s> role to perform this command.", gs.AdminRole), ctx)
}

func SendPermissionsError(gs GuildSettings, ctx *Context) {
	SendError(fmt.Sprintf("You require the <@&%s> role to perform this command.", gs.AccessRole), ctx)
}

func SendBannedError(ctx *Context) {
	SendError("You cannot use this command as you have been banned from using the server", ctx)
}

func SendWrongGuildError(ctx *Context) {
	SendError("You cannot use this bot from outside of the CompSoc server.", ctx)
}

func SendInternalError(err error, ctx *Context) {
	log.Print(err)
	SendError(fmt.Sprintf("An internal error has occurred:\n```\n%s\n```", err), ctx)
}

func CheckGuild(ctx *Context) error {
	guildid := ctx.interaction.GuildId
	if guildid != COMPSOC_GUILD_ID {
		SendWrongGuildError(ctx)
		log.Printf("Guild %s is not the guild (%s).", guildid, os.Getenv("DISCORD_GUILD_ID"))
		return errors.New("Wrong guild.")
	}
	return nil
}

// Vile linear search (grim)
func Contains(arr []string, key string) bool {
	for _, item := range arr {
		if item == key {
			return true
		}
	}
	return false
}

func UserIsAdmin(gs GuildSettings, user *discord.GuildMember) bool {
	return Contains(user.Roles, gs.AdminRole)
}
