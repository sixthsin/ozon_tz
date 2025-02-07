package storage

import (
	"database/sql"
)

type PostgresStorage struct {
	Storage
	db *sql.DB
}

func NewStoragePostgres(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

// func (s *PostgresStorage) GetPosts() ([]*models.Post, error) {
// 	// Реализация SQL-запроса для получения постов
// }

// func (s *PostgresStorage) GetPostByID(id string) (*models.Post, error) {
// 	// Реализация SQL-запроса для получения поста по ID
// }

// func (s *PostgresStorage) CreatePost(post *models.Post) (*models.Post, error) {
// 	// Реализация SQL-запроса для создания поста
// }

// func (s *PostgresStorage) AddComment(comment *models.Comment) (*models.Comment, error) {
// 	// Реализация SQL-запроса для добавления комментария
// }

// func (s *PostgresStorage) GetComments(postID string, first *int, after *string) ([]*models.Comment, error) {
// 	// Реализация SQL-запроса для пагинации комментариев
// }

// func (s *PostgresStorage) GetLatestComment(postID string) (*models.Comment, error) {
// 	// Реализация SQL-запроса для получения последнего комментария
// }
