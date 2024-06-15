package entity

import (
	"testing"
	"github.com/stretchr/testify/assert"
  )

func TestNewUser(t *testing.T) {
	user, err := NewUser("John Doe", "jhon.doe@gmail.com", "password")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.NotEmpty(t, user.Password)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "jhon.doe@gmail.com", user.Email)
}

func TestValidatePassword(t *testing.T) {
	user, err := NewUser("John Doe", "jhon.doe@gmail.com", "password")

	assert.Nil(t, err)
	assert.True(t, user.ValidatePassword("password"))
	assert.False(t, user.ValidatePassword("wrong_password"))
	assert.NotEqual(t, "password", user.Password)
}