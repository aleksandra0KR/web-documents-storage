package handler

import (
	"astral/internal/contorller/utils"
	"astral/internal/model"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

func (h *Handler) UploadDocument(c *gin.Context) {
	metaStr := c.PostForm("meta")
	var meta model.Meta
	var fileData []byte
	var jsonDataBytes []byte
	var err error

	if err := json.Unmarshal([]byte(metaStr), &meta); err != nil {
		c.JSON(model.BadRequestStatusResponse, gin.H{"error": gin.H{"code": model.BadRequestStatusResponse, "text": fmt.Sprintf("invalid meta format: %v", err)}})
		return
	}
	if meta.File == true {
		file, err := c.FormFile("file")

		if err != nil {
			c.JSON(model.BadRequestStatusResponse, gin.H{"error": gin.H{"code": model.BadRequestStatusResponse, "text": fmt.Sprintf("file not provided or invalid: %v", err)}})
			return
		}

		openedFile, err := file.Open()
		if err != nil {
			c.JSON(model.InternalServerErrorStatusResponse, gin.H{"error": gin.H{"code": model.InternalServerErrorStatusResponse, "text": fmt.Sprintf("failed to open file: %v", err)}})
			return
		}
		defer func(openedFile multipart.File) {
			err := openedFile.Close()
			if err != nil {
				log.Errorf("failed to close file: %v", err)
			}
		}(openedFile)

		fileData, err = io.ReadAll(openedFile)
		if err != nil {
			c.JSON(model.InternalServerErrorStatusResponse, gin.H{"error": gin.H{"code": model.InternalServerErrorStatusResponse, "text": fmt.Sprintf("failed to read file: %v", err)}})
			return
		}
	}

	jsonData := make(map[string]interface{})
	jsonStr := c.PostForm("json")
	if jsonStr != "" {
		if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
			c.JSON(model.BadRequestStatusResponse, gin.H{"error": gin.H{"code": model.BadRequestStatusResponse, "text": fmt.Sprintf("invalid json format: %v", err)}})
			return
		}
		jsonDataBytes, err = json.Marshal(jsonData)
		if err != nil {
			c.JSON(model.InternalServerErrorStatusResponse, gin.H{"error": gin.H{"code": model.InternalServerErrorStatusResponse, "text": fmt.Sprintf("failed to upload document: %v", err)}})
			return
		}
	}

	grantData, _ := json.Marshal(meta.Grant)
	doc := model.Document{
		Name:     meta.Name,
		FileData: fileData,
		Public:   meta.Public,
		Mime:     meta.Mime,
		File:     meta.File,
		Token:    meta.Token,
		Grant:    string(grantData),
		JSONData: string(jsonDataBytes),
	}

	err = h.service.UploadDocument(&doc)
	if err != nil {
		return
	}

	if string(jsonDataBytes) != "" && fileData != nil {
		c.JSON(model.SuccessfulStatusResponse, gin.H{"response": gin.H{"json": jsonData, "file": meta.Name}})
	} else if fileData == nil {
		c.JSON(model.SuccessfulStatusResponse, gin.H{"response": gin.H{"json": jsonData}})
	} else {
		c.JSON(model.SuccessfulStatusResponse, gin.H{"data": gin.H{"file": meta.Name}})
	}
}

func (h *Handler) GetDocumentByID(c *gin.Context) {
	id := c.Param("id")
	cookie, err := c.Cookie("token")
	if err != nil {
		c.JSON(model.UnauthorizedStatusResponse, gin.H{"error": gin.H{"code": model.UnauthorizedStatusResponse, "text": "unauthorized"}})
		return
	}
	claims, err := utils.ParseToken(cookie)
	if err != nil {
		c.JSON(model.UnauthorizedStatusResponse, gin.H{"error": gin.H{"code": model.UnauthorizedStatusResponse, "text": err.Error()}})
		return
	}
	document, err := h.service.GetDocumentByID(id, claims.Login)
	if err != nil {
		if string(err.Error()) == "invalid grant" {
			c.JSON(model.ForbiddenStatusResponse, gin.H{"error": gin.H{"code": model.ForbiddenStatusResponse, "text": err.Error()}})
			return
		}
		c.JSON(model.InternalServerErrorStatusResponse, gin.H{"error": gin.H{"code": model.InternalServerErrorStatusResponse, "text": err.Error()}})
		return
	}
	if c.Request.Method == http.MethodHead {
		c.Header("Content-Type", document.Mime)
		c.Header("Content-Length", strconv.Itoa(len(document.FileData)))
		c.Status(model.SuccessfulStatusResponse)
		return
	}
	if document.FileData != nil {
		tempFile, err := os.CreateTemp("", "uploaded-*.")
		if err != nil {
			c.JSON(model.InternalServerErrorStatusResponse, gin.H{"error": gin.H{"code": model.InternalServerErrorStatusResponse, "text": err.Error()}})
			return
		}
		defer func(name string) {
			err := os.Remove(name)
			if err != nil {
				log.Error(err)
			}
		}(tempFile.Name())

		if _, err := tempFile.Write(document.FileData); err != nil {
			c.JSON(model.InternalServerErrorStatusResponse, gin.H{"error": gin.H{"code": model.InternalServerErrorStatusResponse, "text": err.Error()}})
			return
		}
		err = tempFile.Close()
		if err != nil {
			log.Error(err)
		}
		c.FileAttachment(tempFile.Name(), document.Name)
		return
	}

	var jsonData map[string]interface{}
	err = json.Unmarshal([]byte(document.JSONData), &jsonData)
	if err != nil {
		c.JSON(model.InternalServerErrorStatusResponse, gin.H{"error": gin.H{"code": model.InternalServerErrorStatusResponse, "text": err.Error()}})
		return
	}
	c.JSON(model.SuccessfulStatusResponse, gin.H{"response": jsonData})
}

func (h *Handler) GetDocuments(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		c.JSON(model.UnauthorizedStatusResponse, gin.H{"error": gin.H{"code": model.UnauthorizedStatusResponse, "text": "unauthorized"}})
		return
	}
	claims, err := utils.ParseToken(cookie)
	if err != nil {
		c.JSON(model.UnauthorizedStatusResponse, gin.H{"error": gin.H{"code": model.UnauthorizedStatusResponse, "text": err.Error()}})
		return
	}

	login := c.Query("login")
	if login == "" {
		login = claims.Login
	}
	key := c.Query("key")
	value := c.Query("value")
	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	documents, err := h.service.DocumentUsecase.GetDocuments(login, key, value, limit)
	if err != nil {
		c.JSON(model.InternalServerErrorStatusResponse, gin.H{"error": gin.H{"code": model.InternalServerErrorStatusResponse, "text": err.Error()}})
		return
	}

	var responseDocs []gin.H
	var size int
	for _, doc := range documents {
		var grant []string
		err = json.Unmarshal([]byte(doc.Grant), &grant)
		size += len(responseDocs)
		responseDocs = append(responseDocs, gin.H{
			"id":      doc.ID,
			"name":    doc.Name,
			"mime":    doc.Mime,
			"file":    doc.File,
			"public":  doc.Public,
			"created": doc.CreatedAt,
			"grant":   grant,
		})
	}

	response := gin.H{
		"data": gin.H{
			"docs": responseDocs,
		},
	}

	if c.Request.Method == http.MethodHead {
		c.Header("Content-Type", "application/json")
		c.Header("Content-Length", strconv.Itoa(size))
		c.Status(model.SuccessfulStatusResponse)
		return
	}

	c.JSON(model.SuccessfulStatusResponse, gin.H{"response": response})
}

func (h *Handler) DeleteDocumentByID(c *gin.Context) {
	id := c.Param("id")
	cookie, err := c.Cookie("token")
	if err != nil {

		c.JSON(model.UnauthorizedStatusResponse, gin.H{"error": gin.H{"code": model.UnauthorizedStatusResponse, "text": "unauthorized"}})
		return
	}
	claims, err := utils.ParseToken(cookie)
	if err != nil {
		c.JSON(model.UnauthorizedStatusResponse, gin.H{"error": gin.H{"code": model.UnauthorizedStatusResponse, "text": err.Error()}})
	}
	err, token := h.service.DeleteDocumentByID(id, claims.Login)

	if err != nil {
		c.JSON(model.InternalServerErrorStatusResponse, gin.H{"error": gin.H{"code": model.InternalServerErrorStatusResponse, "text": err.Error()}})
		return
	}

	c.JSON(model.SuccessfulStatusResponse, gin.H{"response": gin.H{token: true}})
}
