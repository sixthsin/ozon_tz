package storage

import (
	"context"
	"errors"
	"ozontz/app/models"
	"strconv"
	"sync"
	"time"
)

const textLen = 2000

type InMemoryStorage struct {
	Storage
	mu        sync.Mutex
	posts     map[string]*models.Post
	comments  map[string]*models.Comment
	idCounter int
}

func generateID(contentType string, counter int) string {
	return contentType + strconv.Itoa(counter)
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
	post.ID = generateID("post-", s.idCounter)
	post.CreatedAt = time.Now()
	s.posts[post.ID] = post
	return post, nil
}

func (s *InMemoryStorage) AddComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(comment.Text) > textLen {
		return nil, errors.New("comment text exceeds 2000 characters")
	}

	if _, exists := s.posts[comment.PostID]; !exists {
		return nil, errors.New("post not found")
	}

	s.idCounter++
	comment.ID = generateID("com-", s.idCounter)
	comment.CreatedAt = time.Now()
	s.comments[comment.ID] = comment
	return comment, nil
}

func (s *InMemoryStorage) GetLatestComment(ctx context.Context, postId string) (*models.Comment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var lastComment *models.Comment
	for _, comment := range s.comments {
		if comment.PostID == postId {
			if lastComment == nil || comment.CreatedAt.After(lastComment.CreatedAt) {
				lastComment = comment
			}
		}
	}
	return lastComment, nil
}

// func (s *InMemoryStorage) GetComments(postID string, first *int, after *string) ([]*models.Comment, error) {
// 	// Реализация пагинации комментариев
// }

// func (s *InMemoryStorage) GetLatestComment(postID string) (*models.Comment, error) {
// 	// Реализация получения последнего комментария
// }
