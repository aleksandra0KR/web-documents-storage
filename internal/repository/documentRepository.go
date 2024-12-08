package repository

import "astral/internal/model"

type DocumentRepository interface {
	UploadDocument(*model.Document) error
	GetDocumentByID(string) (*model.Document, error)
	GetDocuments(string, string, string, int) ([]model.Document, error)
	DeleteDocumentByID(string, string) (error, string)
}
