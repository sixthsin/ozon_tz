package storage

import (
	"context"
	"errors"
	"ozontz/app/models"
	"strconv"
	"sync"
	"time"
)

type InMemoryStorage struct {
	Storage
	mu        sync.Mutex
	posts     map[string]*models.Post
	comments  map[string]*models.Comment
	idCounter int
}

func generateID(counter int) string {
	return "id-" + strconv.Itoa(counter)
}

func NewStorageInMemory() *InMemoryStorage {
	return &InMemoryStorage{
		posts:    make(map[string]*models.Post),
		comments: make(map[string]*models.Comment),
	}
}

func (s *InMemoryStorage) GetPosts(ctx context.Context) ([]*models.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var posts []*models.Post
	for _, post := range s.posts {
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *InMemoryStorage) GetPostByID(ctx context.Context, id string) (*models.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post, exists := s.posts[id]
	if !exists {
		return nil, errors.New("post not found")
	}
	return post, nil
}

func (s *InMemoryStorage) CreatePost(ctx context.Context, post *models.Post) (*models.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.idCounter++
	post.ID = generateID(s.idCounter)
	post.CreatedAt = time.Now()
	s.posts[post.ID] = post
	return post, nil
}

// func (s *InMemoryStorage) AddComment(comment *models.Comment) (*models.Comment, error) {
// 	// Реализация добавления комментария
// }

// func (s *InMemoryStorage) GetComments(postID string, first *int, after *string) ([]*models.Comment, error) {
// 	// Реализация пагинации комментариев
// }

// func (s *InMemoryStorage) GetLatestComment(postID string) (*models.Comment, error) {
// 	// Реализация получения последнего комментария
// }
