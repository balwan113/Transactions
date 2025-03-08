package repository

import (
	"context"
	"errors"
	"log"
	"petproject2/cmd/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.Users) (int, error)
	GetUserByName(ctx context.Context, name string) (*models.Users, error)
}

type RepositoryUser struct {
	db *pgxpool.Pool
}

func NewRepositoryUser(db *pgxpool.Pool) *RepositoryUser {
	return &RepositoryUser{db: db}
}

func (r *RepositoryUser) CreateUser(ctx context.Context, users *models.Users) (int, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		log.Printf("Ошибка с Транзакцией: %v", err)
		return 0, err
	}
	defer tx.Rollback(ctx)

	var userId int
	query := `INSERT INTO users (name, password, balance) VALUES ($1, $2, $3) RETURNING id`

	if err := tx.QueryRow(ctx, query, users.Name, users.Password, users.Balance).Scan(&userId); err != nil {
		log.Printf("Ошибка с Записью в БД: %v", err)
		return 0, err
	}
	if err := tx.Commit(ctx); err != nil {
		log.Printf("Ошибка с Коммитом: %v", err)
		return 0, err
	}
	return userId, nil
}

func (r *RepositoryUser) GetUserByName(ctx context.Context, name string) (*models.Users, error) {
	var user models.Users
	query := `SELECT id, name, password, balance FROM users WHERE name = $1`
	if err := r.db.QueryRow(ctx, query, name).Scan(&user.ID, &user.Name, &user.Password, &user.Balance); err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("пользователь не найден")
		}
		return nil, err
	}
	return &user, nil
}
