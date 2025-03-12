package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
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

	userID, ok := atUserID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки user_id"})
		return
	}

	err := h.service.Transfer(ctx, userID, transaction.ToUserID, transaction.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка перевода"})
		return
	}

	c.JSON(200, gin.H{
		"succes": transaction,
	})
}

func (h *HandlerTransations) GetLastTransactions(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неавторизованный пользователь"})
		return
	}

	userID, ok := userIDRaw.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат user_id"})
		return
	}

	log.Printf("Получаем последние транзакции для userID: %d", userID)

	transactions, err := h.service.GetLastTransactions(context.Background(), userID)
	if err != nil {
		log.Printf("Ошибка получения транзакций: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении транзакций"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
