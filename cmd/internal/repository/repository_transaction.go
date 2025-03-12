package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"petproject2/cmd/models"
	"time"
)

type InterfaceTransactions interface {
	TransferMoney(ctx context.Context, transactions *models.Trasactions) error
	GetLastTransactions(ctx context.Context, userId int) ([]models.Trasactions, error)
}

type RepositoryTransactions struct {
	db *pgxpool.Pool
}

func NewTransactionsRepo(db *pgxpool.Pool) *RepositoryTransactions {
	return &RepositoryTransactions{db: db}
}

func (repo *RepositoryTransactions) TransferMoney(ctx context.Context, userid int, to_userid int, amount float64) error {
	tx, err := repo.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		log.Printf("Ошибка с началом транзакции: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	var senderBalance float64
	log.Printf("Перед запросом: ищем баланс пользователя (ID: %d)", userid)
	err = tx.QueryRow(ctx, `SELECT balance FROM users WHERE id = $1 FOR UPDATE`, userid).
		Scan(&senderBalance)
	if err != nil {
		log.Printf("Ошибка при получении баланса пользователя (ID: %d): %v", userid, err)
		return err
	}

	log.Printf("Баланс найден: %f", senderBalance)

	if senderBalance < amount {
		return fmt.Errorf("Ошибка: недостаточно средств")
	}

	_, err = tx.Exec(ctx, `UPDATE users SET balance = balance - $1 WHERE id = $2`, amount, userid)
	if err != nil {
		return fmt.Errorf("Ошибка при списании средств: %v", err)
	}

	_, err = tx.Exec(ctx, `UPDATE users SET balance = balance + $1 WHERE id = $2`, amount, to_userid)
	if err != nil {
		return fmt.Errorf("Ошибка при зачислении денег: %v", err)
	}

	_, err = tx.Exec(ctx, `INSERT INTO transactions(to_userid, at_userid, amount, created_at) VALUES ($1, $2, $3, $4)`,
		to_userid, userid, amount, time.Now())
	if err != nil {
		return fmt.Errorf("Ошибка при записи в таблицу транзакций: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("Ошибка при коммите транзакции: %v", err)
	}

	return err
}
func (repo *RepositoryTransactions) GetLastTransactions(ctx context.Context, userId int) ([]models.Trasactions, error) {
	rows, err := repo.db.Query(ctx, `
SELECT to_userid, at_userid, amount, created_at
	FROM transactions 
	WHERE at_userid = $1 OR to_userid = $1
	ORDER BY created_at DESC  
	LIMIT 10
	`, userId)

	if err != nil {
		return nil, fmt.Errorf("Ошибка при получении транзакций: %v", err)
	}
	defer rows.Close()

	var transactions []models.Trasactions

	for rows.Next() {
		var t models.Trasactions
		if err := rows.Scan(&t.ToUserID, &t.AtUserID, &t.Amount, &t.CreatedAt); err != nil {
			log.Printf("Ошибка при сканировании транзакции: %v", err)
			return nil, fmt.Errorf("Ошибка при сканировании транзакции: %v", err)
		}
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Ошибка при итерации по результатам: %v", err)
	}

	return transactions, nil
}
