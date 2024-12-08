package postgres

import (
	"astral/internal/model"
	"astral/pkg/chekers"
	"fmt"
	"gorm.io/gorm"
)

type DocumentRepositoryPostgres struct {
	db *gorm.DB
}

func NewDocumentRepositoryPostgres(db *gorm.DB) *DocumentRepositoryPostgres {
	return &DocumentRepositoryPostgres{db: db}
}

func (documentRepository *DocumentRepositoryPostgres) UploadDocument(document *model.Document) error {
	return documentRepository.db.Create(document).Error
}

func (documentRepository *DocumentRepositoryPostgres) GetDocumentByID(id string) (*model.Document, error) {
	var doc model.Document
	if err := documentRepository.db.First(&doc, id).Error; err != nil {
		return nil, err
	}
	return &doc, nil
}

func (documentRepository *DocumentRepositoryPostgres) GetDocuments(login, key, value string, limit int) ([]model.Document, error) {
	query := documentRepository.db.Model(&model.Document{})
	if login != "" {
		query = query.Where("documents.grant LIKE ?", "%"+login+"%")
	}
	if key != "" && value != "" {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	var documents []model.Document
	if err := query.Order("name ASC, created_at ASC").Limit(limit).Find(&documents).Error; err != nil {
		return nil, err
	}
	return documents, nil
}

func (documentRepository *DocumentRepositoryPostgres) DeleteDocumentByID(id, login string) (error, string) {
	var document model.Document
	err := documentRepository.db.Select("public").First(&document, id).Error
	if err != nil {
		return err, ""
	}

	if !document.Public {
		err = documentRepository.db.Select("grant").First(&document, id).Error
		if err != nil {
			return err, ""
		}
		valid, err := chekers.ContainsInString(document.Grant, login)
		if err != nil {
			return err, ""
		}
		if !valid {
			return fmt.Errorf("invalid grant"), ""
		}
	}

	err = documentRepository.db.Select("token").First(&document, id).Error
	if err != nil {
		return err, ""
	}
	return documentRepository.db.Delete(&document, id).Error, document.Token
}
