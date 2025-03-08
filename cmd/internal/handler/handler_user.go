package handler

import (
	"context"
	"net/http"
	"petproject2/cmd/internal/service"
	"petproject2/cmd/models"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	userService service.UserService
}

func NewHandler(userService service.UserService) *Handler {
	return &Handler{userService: userService}
}

func (h *Handler) RegisterUser(c *gin.Context) {
	var user models.Users

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}

	id, err := h.userService.RegisterUser(context.Background(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при регистрации"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Регистрация успешна", "user_id": id})
}

func (h *Handler) LoginUser(c *gin.Context) {
	var user models.Users

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}

	token, err := h.userService.LoginUser(context.Background(), user.Name, user.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль "})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
