package storage

import (
	"context"
	"database/sql"
	"ozontz/app/models"
)

type Storage interface {
	NewStorageInMemory() *InMemoryStorage
	NewStoragePostgres(db *sql.DB) *PostgresStorage
	GetPosts(ctx context.Context) ([]*models.Post, error)
	GetPostByID(ctx context.Context, id string) (*models.Post, error)
	CreatePost(ctx context.Context, post *models.Post) (*models.Post, error)
	AddComment(ctx context.Context, comment *models.Comment) (*models.Comment, error)
	GetComments(ctx context.Context, postId string, first *int, after *string) ([]*models.Comment, error)
	GetLatestComment(ctx context.Context, postId string) (*models.Comment, error)
}
