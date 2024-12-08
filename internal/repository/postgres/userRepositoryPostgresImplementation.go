package postgres

import (
	"astral/internal/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type UserRepositoryPostgres struct {
	db *gorm.DB
}

func NewUserRepositoryPostgres(db *gorm.DB) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{db: db}
}

func (userRepository *UserRepositoryPostgres) RegisterUser(user *model.User) error {
	var existingUser model.User
	userRepository.db.Where("login = ?", user.Login).Find(&existingUser)
	if existingUser.Login != "" {
		return fmt.Errorf("user already exists")
	}

	userRepository.db.Create(&user)
	return nil
}

func (userRepository *UserRepositoryPostgres) Auth(user *model.User) (error, *model.User) {
	var existingUser model.User
	userRepository.db.Where("login = ?", user.Login).Find(&existingUser)
	if existingUser.Login == "" {
		return errors.New("user doesn't exist"), nil
	}
	return nil, &existingUser
}
