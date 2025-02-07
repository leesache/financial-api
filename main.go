// main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/leesache/financial-api/handler"
	"github.com/leesache/financial-api/repository"
	"github.com/leesache/financial-api/service"
)

func waitForDB(ctx context.Context, connString string) (*pgx.Conn, error) {
	var db *pgx.Conn
	var err error

	retries := 10
	for i := 0; i < retries; i++ {
		db, err = pgx.Connect(ctx, connString)
		if err == nil {
			log.Println("Successfully connected to the database!")
			return db, nil
		}
		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, retries, err)
		time.Sleep(3 * time.Second) // Wait before retrying
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts", retries)
}

func main() {
	// Get database connection details from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Wait for the database to become available
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := waitForDB(ctx, connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	// Initialize Repository
	repo := repository.NewPGAccountRepository(conn)

	// Initialize Service
	svc := service.NewAccountService(repo)

	// Initialize Handler
	accountHandler := handler.NewAccountHandler(svc)

	// Initialize Gin Router
	router := gin.Default()

	// Define API routes
	api := router.Group("/api")
	{
		api.GET("/account/:id", accountHandler.GetAccount)
		api.POST("/transfer", accountHandler.TransferFunds)
		api.GET("/account/:id/transactions", accountHandler.GetTransactionHistory)
	}

	// Start the HTTP server
	log.Println("Starting server on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
