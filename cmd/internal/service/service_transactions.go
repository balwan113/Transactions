package service

import (
	"context"
	"fmt"
	"petproject2/cmd/internal/repository"
	"petproject2/cmd/models"
)

type InterfaceServiceTransactions interface {
	Transfer(ctx context.Context, transactions *models.Trasactions) error
}
type ServiceTransactions struct {
	s *repository.RepositoryTransactions
}

func NewServiceTransactions(s *repository.RepositoryTransactions) *ServiceTransactions {
	return &ServiceTransactions{s: s}
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
