package database

import (
	"github.com/drawiin/go-expert/crud/internal/entity"
	"gorm.io/gorm"
)

type UserDBInterface interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
}

type UserDB struct {
	DB *gorm.DB
}

func NewUserDB(db *gorm.DB) *UserDB {
	return &UserDB{
		DB: db,
	}
}

func (db *UserDB) Create(user *entity.User) error {
	return db.DB.Create(user).Error
}

func (db *UserDB) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
