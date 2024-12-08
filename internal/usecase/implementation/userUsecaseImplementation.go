package implementation

import (
	"astral/internal/contorller/utils"
	"astral/internal/model"
	"astral/internal/repository"
	"astral/pkg/chekers"
	"errors"
	"fmt"
)

type UserUsecaseImplementation struct {
	repository repository.UserRepository
}

func NewUserUsecaseImplementation(repository repository.UserRepository) *UserUsecaseImplementation {
	return &UserUsecaseImplementation{repository: repository}
}

func (userUsecase *UserUsecaseImplementation) RegisterUser(user *model.User) error {
	err := chekers.CheckLoginValidation(user.Login)
	if err != nil {
		return err
	}
	err = chekers.CheckPasswordValidation(user.Password)
	if err != nil {
		return err
	}

	user.Password, err = utils.GenerateHashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("could not generate password hash")
	}
	err = userUsecase.repository.RegisterUser(user)
	return err
}

func (userUsecase *UserUsecaseImplementation) Auth(user *model.User) (error, *model.User) {
	err, existingUser := userUsecase.repository.Auth(user)
	if err != nil {
		return err, nil
	}

	err = utils.CompareHashPassword(user.Password, existingUser.Password)
	if err != nil {
		return errors.New("invalid password"), nil
	}
	return nil, existingUser
}
