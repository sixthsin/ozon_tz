package storage

import (
	"context"
	"errors"
	"ozontz/app/models"
	"sort"
	"strconv"
	"sync"
	"time"
)

const (
	textLen       = 2000
	commentsCount = 5
)

type InMemoryStorage struct {
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

	posts := make([]*models.Post, 5)
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
		if comment.PostID == postId && comment.ParentID == nil {
			if lastComment == nil || comment.CreatedAt.After(lastComment.CreatedAt) {
				lastComment = comment
			}
		}
	}
	return lastComment, nil
}

func (s *InMemoryStorage) GetComments(ctx context.Context, postId string, after *string) ([]*models.Comment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var comments []*models.Comment
	for _, comment := range s.comments {
		if comment.PostID == postId {
			comments = append(comments, comment)
		}
	}

	sort.Slice(comments, func(i, j int) bool {
		return comments[i].CreatedAt.Before(comments[j].CreatedAt)
	})

	index := -1
	if after != nil {

		for i, comment := range comments {
			if generateCursor(comment) == *after {
				index = i
				break
			}
		}
		if index == -1 {
			return nil, errors.New("invalid cursor")
		}
		if index != -1 {
			comments = comments[index+1:]
		}
	} else {
		if len(comments) > commentsCount {
			comments = comments[:commentsCount]
		}
	}

	return comments, nil
}

func generateCursor(comment *models.Comment) string {
	return "cur-" + comment.ID
}
