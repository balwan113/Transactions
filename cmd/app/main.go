package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"log"
	"petproject2/cmd/internal/handler"
	"petproject2/cmd/internal/repository"
	"petproject2/cmd/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func initRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}

	log.Println("Redis успешно подключен!")
	return rdb
}

func main() {
	dbUrl := "postgres://postgres:kilboo123@localhost:5432/testdb?sslmode=disable"

	dbPool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}
	defer dbPool.Close()

	rdb := initRedis()
	defer rdb.Close()

	repo := repository.NewRepositoryUser(dbPool)
	userService := service.NewServiceUser(repo)
	h := handler.NewHandler(userService)

	repoT := repository.NewTransactionsRepo(dbPool)
	transactionsService := service.NewServiceTransactions(repoT, rdb)
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
		auth.GET("/getlast", hs.GetLastTransactions)
	}

	r.Run(":8080")
}
