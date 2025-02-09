package graph

import (
	"context"
	"errors"
	"testing"

	"ozontz/app/models"

	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
)

type MockStorage struct {
	CreatePostFn       func(ctx context.Context, post *models.Post) (*models.Post, error)
	GetPostByIDFn      func(ctx context.Context, id string) (*models.Post, error)
	GetPostsFn         func(ctx context.Context) ([]*models.Post, error)
	AddCommentFn       func(ctx context.Context, comment *models.Comment) (*models.Comment, error)
	GetLatestCommentFn func(ctx context.Context, postId string) (*models.Comment, error)
	GetCommentsFn      func(ctx context.Context, postId string, after *string) ([]*models.Comment, error)
}

func (m *MockStorage) CreatePost(ctx context.Context, post *models.Post) (*models.Post, error) {
	return m.CreatePostFn(ctx, post)
}

func (m *MockStorage) GetPostByID(ctx context.Context, id string) (*models.Post, error) {
	return m.GetPostByIDFn(ctx, id)
}

func (m *MockStorage) GetPosts(ctx context.Context) ([]*models.Post, error) {
	return m.GetPostsFn(ctx)
}

func (m *MockStorage) AddComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	return m.AddCommentFn(ctx, comment)
}

func (m *MockStorage) GetLatestComment(ctx context.Context, postId string) (*models.Comment, error) {
	return m.GetLatestCommentFn(ctx, postId)
}

func (m *MockStorage) GetComments(ctx context.Context, postId string, after *string) ([]*models.Comment, error) {
	return m.GetCommentsFn(ctx, postId, after)
}

func TestResolveCreatePost(t *testing.T) {
	mockStore := &MockStorage{
		CreatePostFn: func(ctx context.Context, post *models.Post) (*models.Post, error) {
			return post, nil
		},
	}
	SetStore(mockStore)

	params := graphql.ResolveParams{
		Args: map[string]interface{}{
			"title":         "Test Post",
			"content":       "This is a test post.",
			"authorId":      "user-1",
			"allowComments": true,
		},
	}

	result, err := resolveCreatePost(params)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	post, ok := result.(*models.Post)
	assert.True(t, ok)
	assert.Equal(t, "Test Post", post.Title)
	assert.Equal(t, "This is a test post.", post.Content)
	assert.Equal(t, "user-1", post.AuthorID)
	assert.True(t, post.AllowComments)
}

func TestResolveGetPost(t *testing.T) {
	mockStore := &MockStorage{
		GetPostByIDFn: func(ctx context.Context, id string) (*models.Post, error) {
			if id == "post-1" {
				return &models.Post{ID: "post-1", Title: "Test Post"}, nil
			}
			return nil, errors.New("post not found")
		},
	}
	SetStore(mockStore)

	t.Run("Valid ID", func(t *testing.T) {
		params := graphql.ResolveParams{
			Args: map[string]interface{}{
				"id": "post-1",
			},
		}

		result, err := resolveGetPost(params)
		assert.NoError(t, err)
		assert.NotNil(t, result)

		post, ok := result.(*models.Post)
		assert.True(t, ok)
		assert.Equal(t, "post-1", post.ID)
		assert.Equal(t, "Test Post", post.Title)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		params := graphql.ResolveParams{
			Args: map[string]interface{}{
				"id": "invalid-id",
			},
		}

		result, err := resolveGetPost(params)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestResolveGetPostsList(t *testing.T) {
	mockStore := &MockStorage{
		GetPostsFn: func(ctx context.Context) ([]*models.Post, error) {
			return []*models.Post{
				{ID: "post-1", Title: "Post 1"},
				{ID: "post-2", Title: "Post 2"},
			}, nil
		},
	}
	SetStore(mockStore)

	params := graphql.ResolveParams{}

	result, err := resolveGetPostsList(params)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	posts, ok := result.([]*models.Post)
	assert.True(t, ok)
	assert.Len(t, posts, 2)
	assert.Equal(t, "Post 1", posts[0].Title)
	assert.Equal(t, "Post 2", posts[1].Title)
}

func TestResolveAddComment(t *testing.T) {
	mockStore := &MockStorage{
		AddCommentFn: func(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
			return comment, nil
		},
	}
	SetStore(mockStore)

	params := graphql.ResolveParams{
		Args: map[string]interface{}{
			"postId":   "post-1",
			"parentId": nil,
			"authorId": "user-1",
			"text":     "Test comment",
		},
	}

	result, err := resolveAddComment(params)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	comment, ok := result.(*models.Comment)
	assert.True(t, ok)
	assert.Equal(t, "post-1", comment.PostID)
	assert.Nil(t, comment.ParentID)
	assert.Equal(t, "user-1", comment.AuthorID)
	assert.Equal(t, "Test comment", comment.Text)
}

func TestResolveGetLastComment(t *testing.T) {
	mockStore := &MockStorage{
		GetLatestCommentFn: func(ctx context.Context, postId string) (*models.Comment, error) {
			if postId == "post-1" {
				return &models.Comment{ID: "com-1", Text: "Latest comment"}, nil
			}
			return nil, errors.New("no comments found")
		},
	}
	SetStore(mockStore)

	t.Run("Valid Post", func(t *testing.T) {
		post := &models.Post{ID: "post-1", AllowComments: true}
		params := graphql.ResolveParams{
			Source: post,
		}

		result, err := resolveGetLastComment(params)
		assert.NoError(t, err)
		assert.NotNil(t, result)

		comment, ok := result.(*models.Comment)
		assert.True(t, ok)
		assert.Equal(t, "com-1", comment.ID)
		assert.Equal(t, "Latest comment", comment.Text)
	})

	t.Run("No Comments", func(t *testing.T) {
		post := &models.Post{ID: "post-2", AllowComments: true}
		params := graphql.ResolveParams{
			Source: post,
		}

		result, err := resolveGetLastComment(params)
		assert.Nil(t, result)
		if err == nil || err.Error() != "no comments found" {
			t.Errorf("Expected error 'no comments found', but got: %v", err)
		}
	})
}

func TestResolveGetComments(t *testing.T) {
	mockStore := &MockStorage{
		GetCommentsFn: func(ctx context.Context, postId string, after *string) ([]*models.Comment, error) {
			if postId == "post-1" {
				return []*models.Comment{
					{ID: "com-1", Text: "Comment 1"},
					{ID: "com-2", Text: "Comment 2"},
				}, nil
			}
			return nil, errors.New("no comments found")
		},
	}
	SetStore(mockStore)

	t.Run("Valid Post ID", func(t *testing.T) {
		params := graphql.ResolveParams{
			Args: map[string]interface{}{
				"postId": "post-1",
			},
		}

		result, err := resolveGetComments(params)
		assert.NoError(t, err)
		assert.NotNil(t, result)

		comments, ok := result.([]*models.Comment)
		assert.True(t, ok)
		assert.Len(t, comments, 2)
		assert.Equal(t, "Comment 1", comments[0].Text)
		assert.Equal(t, "Comment 2", comments[1].Text)
	})

	t.Run("Invalid Post ID", func(t *testing.T) {
		params := graphql.ResolveParams{
			Args: map[string]interface{}{
				"postId": "invalid-post",
			},
		}

		result, err := resolveGetComments(params)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
