package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"petproject2/cmd/internal/repository"
	"petproject2/cmd/models"
	"strconv"
	"time"
)

type InterfaceServiceTransactions interface {
	Transfer(ctx context.Context, transactions *models.Trasactions) error
}
type ServiceTransactions struct {
	s           *repository.RepositoryTransactions
	redisClient *redis.Client
}

func NewServiceTransactions(s *repository.RepositoryTransactions, redisClient *redis.Client) *ServiceTransactions {
	return &ServiceTransactions{s: s, redisClient: redisClient}
}

func (s *ServiceTransactions) Transfer(ctx context.Context, userid int, to_userid int, amount float64) error {
	if amount < 0 {
		return fmt.Errorf("Сумма не Должна быть отрицательной!")
	}
	if userid == to_userid {
		return fmt.Errorf("Нельзя перевести деньги самому себе!")
	}
	return s.s.TransferMoney(ctx, userid, to_userid, amount)
}

func (s *ServiceTransactions) GetLastTransactions(ctx context.Context, userID int) ([]models.Trasactions, error) {
	cachedKey := strconv.Itoa(userID)

	cachedData, err := s.redisClient.Get(ctx, cachedKey).Result()
	if err == nil {
		var transactions []models.Trasactions
		if json.Unmarshal([]byte(cachedData), &transactions) == nil {
			fmt.Println("Редис Записан")
			return transactions, nil
		}

	}

	transactions, err := s.s.GetLastTransactions(ctx, userID)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(transactions)
	err = s.redisClient.Set(ctx, cachedKey, data, 10*time.Minute).Err()
	if err != nil {
		log.Printf("Ошибка записи в Redis: %v", err)
	}

	return transactions, nil
}
