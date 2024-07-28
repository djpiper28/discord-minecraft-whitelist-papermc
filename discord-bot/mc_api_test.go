package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMcApiSuccess(t *testing.T) {
  user, err := GetMinecraftUser("DifficultPete")
  assert.NoError(t, err)
  assert.Equal(t, "DifficultPete", user.Name)
  assert.NotEmpty(t, user.Id)
}
