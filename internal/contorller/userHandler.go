package handler

import (
	"astral/internal/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

func (h *Handler) Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(model.BadRequestStatusResponse, gin.H{"error": gin.H{"code": model.BadRequestStatusResponse, "text": "invalid input"}})
		return
	}

	if user.Token != os.Getenv("ADMIN_TOKEN") {
		c.JSON(model.UnauthorizedStatusResponse, gin.H{"error": gin.H{"code": model.UnauthorizedStatusResponse, "text": "unauthorized error"}})
		return
	}

	if err := h.service.UserUsecase.RegisterUser(&user); err != nil {
		c.JSON(model.BadRequestStatusResponse, gin.H{"error": gin.H{"code": model.BadRequestStatusResponse, "text": err.Error()}})
		return
	}
	c.JSON(model.SuccessfulStatusResponse, gin.H{"response": gin.H{"login": user.Login}})
}

func (h *Handler) Auth(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(model.BadRequestStatusResponse, gin.H{"error": gin.H{"code": model.BadRequestStatusResponse, "text": "invalid input"}})
		return
	}

	err, existingUser := h.service.UserUsecase.Auth(&user)
	if err != nil {
		c.JSON(model.UnauthorizedStatusResponse, gin.H{"error": gin.H{"code": model.UnauthorizedStatusResponse, "text": err.Error()}})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &model.Claims{
		Login: existingUser.Login,
		StandardClaims: jwt.StandardClaims{
			Subject:   existingUser.Login,
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		c.JSON(model.InternalServerErrorStatusResponse, gin.H{"error": gin.H{"code": model.InternalServerErrorStatusResponse, "text": "could not generate token"}})
		return
	}
	c.SetCookie("token", tokenString, int(expirationTime.Unix()), "/", "localhost", false, true)
	c.JSON(model.SuccessfulStatusResponse, gin.H{"response": gin.H{"token": tokenString}})
}

func (h *Handler) Logout(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(model.UnauthorizedStatusResponse, gin.H{"error": gin.H{"code": model.UnauthorizedStatusResponse, "text": "unauthorized"}})
		return
	}
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(model.SuccessfulStatusResponse, gin.H{"response": gin.H{token: "true"}})
}
