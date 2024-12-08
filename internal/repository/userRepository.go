package repository

import "astral/internal/model"

type UserRepository interface {
	RegisterUser(user *model.User) error
	Auth(user *model.User) (error, *model.User)
}
