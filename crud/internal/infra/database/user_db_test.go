package database

import (
	"testing"

	"github.com/drawiin/go-expert/crud/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	db := CreateUserTestDB(t)
	user, _ := entity.NewUser("John Doe", "j@j.com", "password")
	userDB := NewUserDB(db)

	err := userDB.Create(user)
	assert.Nil(t, err)

	var userFound entity.User
	err = db.First(&userFound, "id = ?", user.ID).Error

	assert.Nil(t, err)
	assert.NotNil(t, userFound)
	assert.Equal(t, user.Email, userFound.Email)
	assert.Equal(t, user.Name, userFound.Name)
	assert.Equal(t, user.Name, userFound.Name)
	assert.NotEmpty(t, user.Password)
}

func TestFindByEmail(t *testing.T) {
	db := CreateUserTestDB(t)
	user, _ := entity.NewUser("John Doe", "j@j.com", "password")
	userDB := NewUserDB(db)
	userDB.Create(user)

	userFound, err := userDB.FindByEmail("j@j.com")

	assert.Nil(t, err)
	assert.Equal(t, user.Email, userFound.Email)
	assert.Equal(t, user.Name, userFound.Name)
}

func TestFindByEmailNotFound(t *testing.T) {
	db := CreateUserTestDB(t)
	user, _ := entity.NewUser("John Doe", "j@j.com", "password")
	userDB := NewUserDB(db)
	userDB.Create(user)

	userFound, err := userDB.FindByEmail("x@x.com")

	assert.NotNil(t, err)
	assert.Nil(t, userFound)
}

func CreateUserTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error connecting to the database: %v", err)
	}
	db.AutoMigrate(&entity.User{})
	return db
}
