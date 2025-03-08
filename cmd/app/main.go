package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"petproject2/cmd/internal/handler"
	"petproject2/cmd/internal/repository"
	"petproject2/cmd/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Подключение к базе данных
	dbUrl := ""

	dbPool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}
	defer dbPool.Close()

	// Инициализация репозитория, сервиса и обработчика
	repo := repository.NewRepositoryUser(dbPool)
	userService := service.NewServiceUser(repo)
	h := handler.NewHandler(userService)
	repoT := repository.NewTransactionsRepo(dbPool)
	transactionsService := service.NewServiceTransactions(repoT)
	hs := handler.NewHandlerTransactions(transactionsService)

	r := gin.Default()

	userGroup := r.Group("/")
	{
		userGroup.POST("/login", h.LoginUser)
		userGroup.POST("/register", h.RegisterUser)
	}
	auth := r.Group("/user")
	auth.Use(service.AuthMiddleware())
	{
		auth.POST("/transfer", hs.TransferMoney)
	}
	r.Run(":8080")
}
