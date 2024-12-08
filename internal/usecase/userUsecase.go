package usecase

import "astral/internal/model"

type UserUsecase interface {
	RegisterUser(user *model.User) error
	Auth(user *model.User) (error, *model.User)
}
