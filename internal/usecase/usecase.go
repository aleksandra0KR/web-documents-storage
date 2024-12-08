package usecase

import (
	"astral/internal/repository"
	"astral/internal/usecase/implementation"
)

type UseCase struct {
	DocumentUsecase
	UserUsecase
}

func NewUseCase(repository *repository.Repository) *UseCase {
	return &UseCase{
		DocumentUsecase: implementation.NewDocumentUsecaseImplementation(repository.DocumentRepository),
		UserUsecase:     implementation.NewUserUsecaseImplementation(repository.UserRepository),
	}
}
