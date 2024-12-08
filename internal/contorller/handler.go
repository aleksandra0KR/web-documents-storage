package handler

import (
	"astral/internal/model"
	"astral/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	service *usecase.UseCase
}

func NewHandler(service *usecase.UseCase) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Handle() http.Handler {
	router := gin.Default()

	router.POST("/api/register", h.Register)
	router.POST("/api/auth", h.Auth)
	router.DELETE("/api/auth/:token", h.Logout)

	router.POST("/api/docs", h.UploadDocument)
	router.GET("/api/docs/:id", h.GetDocumentByID)
	router.HEAD("/api/docs/:id", h.GetDocumentByID)
	router.GET("/api/docs", h.GetDocuments)
	router.HEAD("/api/docs", h.GetDocuments)
	router.DELETE("/api/docs/:id", h.DeleteDocumentByID)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(model.NotImplementedStatusResponse, gin.H{"code": model.NotImplementedStatusResponse, "error": "not implemented"})
	})
	return router
}
