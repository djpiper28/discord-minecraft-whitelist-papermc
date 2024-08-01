FROM golang:1.21

WORKDIR /app
COPY . .

ENV DISCORD_TOKEN=change_me
ENV DISCORD_GUILD_ID=change_me
ENV DATABASE_URL=change_me
RUN go build

ENTRYPOINT ["/app/minecraft-discord-bot"]
