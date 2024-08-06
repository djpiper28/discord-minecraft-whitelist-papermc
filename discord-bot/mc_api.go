package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

type McUser struct {
	Name string    `json:"name"`
	Id   uuid.UUID `json:"id"`
}

func GetMinecraftUser(username string) (McUser, error) {
	usernameSafe := url.QueryEscape(username)
	uri := fmt.Sprintf("https://api.mojang.com/users/profiles/minecraft/%s", usernameSafe)

	resp, err := http.Get(uri)
	if err != nil {
		return McUser{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return McUser{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return McUser{}, fmt.Errorf("Minecraft API returned status code %d: %s", resp.StatusCode, body)
	}

	var user McUser
	err = json.Unmarshal(body, &user)
	if err != nil {
		return McUser{}, err
	}

	return user, nil
}
