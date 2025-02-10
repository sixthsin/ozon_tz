package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"ozontz/app/models"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewStoragePostgres(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (s *PostgresStorage) GetPosts(ctx context.Context) ([]*models.Post, error) {
	query := `
        SELECT id, title, content, author_id, allow_comments, created_at
        FROM posts
        WHERE true
    `

	args := []interface{}{}

	query += " ORDER BY created_at DESC"
	if postsCount > 0 {
		query += " LIMIT $" + strconv.Itoa(len(args)+1)
		args = append(args, postsCount)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.AllowComments, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (s *PostgresStorage) GetPostByID(ctx context.Context, id string) (*models.Post, error) {
	row := s.db.QueryRowContext(ctx, "SELECT id, title, content, author_id, allow_comments, created_at FROM posts WHERE id = $1", id)

	post := &models.Post{}
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.AllowComments, &post.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	return post, nil
}

func (s *PostgresStorage) CreatePost(ctx context.Context, post *models.Post) (*models.Post, error) {
	query := `
        INSERT INTO posts (id, title, content, author_id, allow_comments, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, title, content, author_id, allow_comments, created_at
    `

	post.ID = generateId("post-")
	post.CreatedAt = time.Now().UTC()

	err := s.db.QueryRowContext(ctx, query, post.ID, post.Title, post.Content, post.AuthorID, post.AllowComments, post.CreatedAt).Scan(
		&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.AllowComments, &post.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostgresStorage) AddComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	if len(comment.Text) > 2000 {
		return nil, errors.New("comment text exceeds 2000 characters")
	}

	query := `
        INSERT INTO comments (id, post_id, parent_id, author_id, text, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, post_id, parent_id, author_id, text, created_at
    `

	comment.ID = generateId("com-")
	comment.CreatedAt = time.Now().UTC()

	var parentId sql.NullString
	if comment.ParentID != nil {
		parentId.String = *comment.ParentID
		parentId.Valid = true
	}

	err := s.db.QueryRowContext(ctx, query, comment.ID, comment.PostID, parentId, comment.AuthorID, comment.Text, comment.CreatedAt).Scan(
		&comment.ID, &comment.PostID, &comment.ParentID, &comment.AuthorID, &comment.Text, &comment.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *PostgresStorage) GetComments(ctx context.Context, postId string, fromID *string) ([]*models.Comment, error) {
	query := `
        SELECT id, post_id, parent_id, author_id, text, created_at
        FROM comments
        WHERE post_id = $1
    `

	args := []interface{}{postId}

	if fromID != nil {
		query += " AND created_at > (SELECT created_at FROM comments WHERE id = $" + strconv.Itoa(len(args)+1) + ")"
		args = append(args, *fromID)
	}

	query += " ORDER BY created_at ASC"
	if commentsCount > 0 {
		query += " LIMIT $" + strconv.Itoa(len(args)+1)
		args = append(args, commentsCount)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		comment := &models.Comment{}
		var parentId sql.NullString
		err := rows.Scan(&comment.ID, &comment.PostID, &parentId, &comment.AuthorID, &comment.Text, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		if parentId.Valid {
			comment.ParentID = &parentId.String
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (s *PostgresStorage) GetLatestComment(ctx context.Context, postId string) (*models.Comment, error) {
	query := `
        SELECT id, post_id, parent_id, author_id, text, created_at
        FROM comments
        WHERE post_id = $1 AND parent_id IS NULL
        ORDER BY created_at DESC
        LIMIT 1
    `

	row := s.db.QueryRowContext(ctx, query, postId)

	comment := &models.Comment{}
	err := row.Scan(&comment.ID, &comment.PostID, &comment.ParentID, &comment.AuthorID, &comment.Text, &comment.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return comment, nil
}

func (s *PostgresStorage) ApplyMigrations(migrationDir string) error {
	migrationDir = filepath.ToSlash(strings.ReplaceAll(migrationDir, ":", "|"))
	if !strings.HasPrefix(migrationDir, "file://") {
		migrationDir = "file://" + migrationDir
	}

	dbDriver, err := postgres.WithInstance(s.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to initialize Postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationDir,
		"postgres",
		dbDriver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}

func generateId(contentType string) string {
	return contentType + strconv.FormatInt(time.Now().UnixNano(), 10)
}
