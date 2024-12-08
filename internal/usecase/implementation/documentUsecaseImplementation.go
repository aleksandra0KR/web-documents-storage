package implementation

import (
	"astral/internal/cache"
	"astral/internal/model"
	"astral/internal/repository"
	"astral/pkg/chekers"
	"encoding/json"
	"fmt"
)

type DocumentUsecaseImplementation struct {
	repository repository.DocumentRepository
}

func NewDocumentUsecaseImplementation(repository repository.DocumentRepository) *DocumentUsecaseImplementation {
	return &DocumentUsecaseImplementation{repository: repository}
}

func (documentUsecase *DocumentUsecaseImplementation) UploadDocument(document *model.Document) error {
	err := documentUsecase.repository.UploadDocument(document)
	if err != nil {
		return err
	}
	cacheKey := "document:" + document.ID
	docBytes, _ := json.Marshal(document)
	cache.Cache.Set(cacheKey, docBytes, 1)
	return nil
}

func (documentUsecase *DocumentUsecaseImplementation) GetDocumentByID(id, login string) (*model.Document, error) {
	cacheKey := "document:" + id
	if item, found := cache.Cache.Get(cacheKey); found {
		var doc model.Document
		if err := json.Unmarshal(item.([]byte), &doc); err != nil {
			return nil, err
		}
		return &doc, nil
	}

	document, err := documentUsecase.repository.GetDocumentByID(id)
	if err != nil {
		return nil, err
	}

	if !document.Public {
		valid, err := chekers.ContainsInString(document.Grant, login)
		if err != nil {
			return nil, err
		} else if !valid {
			return nil, fmt.Errorf("invalid grant")
		}
	}

	docBytes, _ := json.Marshal(document)
	cache.Cache.Set(cacheKey, docBytes, 1)
	return document, nil
}

func (documentUsecase *DocumentUsecaseImplementation) GetDocuments(login string, key string, value string, limit int) ([]model.Document, error) {
	documents, err := documentUsecase.repository.GetDocuments(login, key, value, limit)
	if err != nil {
		return nil, err
	}
	return documents, nil
}

func (documentUsecase *DocumentUsecaseImplementation) DeleteDocumentByID(id, login string) (error, string) {
	err, token := documentUsecase.repository.DeleteDocumentByID(id, login)
	cacheKey := "document:" + id
	cache.Cache.Del(cacheKey)
	return err, token
}
