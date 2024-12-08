package usecase

import "astral/internal/model"

type DocumentUsecase interface {
	UploadDocument(*model.Document) error
	GetDocumentByID(string, string) (*model.Document, error)
	GetDocuments(string, string, string, int) ([]model.Document, error)
	DeleteDocumentByID(string, string) (error, string)
}
