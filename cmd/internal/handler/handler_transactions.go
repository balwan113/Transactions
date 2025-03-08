package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"petproject2/cmd/internal/service"
	"petproject2/cmd/models"
	"time"
)

type HandlerTransations struct {
	service *service.ServiceTransactions
}

func NewHandlerTransactions(service *service.ServiceTransactions) *HandlerTransations {
	return &HandlerTransations{service: service}
}

func (h *HandlerTransations) TransferMoney(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var transaction models.Trasactions
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(400, gin.H{
			"error": "Ошибка с записью в Json",
		})
		return
	}

	atUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Ошибка авторизации"})
		return
	}

	// Конвертируем user_id в int
	userID, ok := atUserID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки user_id"})
		return
	}

	// Дальше передаём userID в сервис для обработки транзакции
	err := h.service.Transfer(ctx, userID, transaction.To_userid, transaction.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка перевода"})
		return
	}

	c.JSON(200, gin.H{
		"succes": transaction,
	})
}
