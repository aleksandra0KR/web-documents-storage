package repository

import (
	"astral/internal/repository/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	DocumentRepository
	UserRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DocumentRepository: postgres.NewDocumentRepositoryPostgres(db),
		UserRepository:     postgres.NewUserRepositoryPostgres(db),
	}
}
