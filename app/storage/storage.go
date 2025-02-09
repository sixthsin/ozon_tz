package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"ozontz/app/models"
)

const (
	textLen       = 2000
	commentsCount = 5
	postsCount    = 10
)

type Storage interface {
	GetPosts(ctx context.Context) ([]*models.Post, error)
	GetPostByID(ctx context.Context, id string) (*models.Post, error)
	CreatePost(ctx context.Context, post *models.Post) (*models.Post, error)
	AddComment(ctx context.Context, comment *models.Comment) (*models.Comment, error)
	GetComments(ctx context.Context, postId string, after *string) ([]*models.Comment, error)
	GetLatestComment(ctx context.Context, postId string) (*models.Comment, error)
}

func InitPostgresDB() (*sql.DB, error) {
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		return nil, fmt.Errorf("environment variable DB_USER is not set")
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		return nil, fmt.Errorf("environment variable DB_PASSWORD is not set")
	}
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		return nil, fmt.Errorf("environment variable DB_HOST is not set")
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		return nil, fmt.Errorf("environment variable DB_PORT is not set")
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		return nil, fmt.Errorf("environment variable DB_NAME is not set")
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	log.Println("Database connection established")

	migrationDir := "./migrations"
	store := NewStoragePostgres(db)
	if err := store.ApplyMigrations(migrationDir); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}
	log.Println("Migrations applied successfully")

	return db, nil
}
