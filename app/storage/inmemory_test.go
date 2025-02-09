package storage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"ozontz/app/models"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryCreatePost(t *testing.T) {
	store := NewStorageInMemory()
	post := &models.Post{
		ID:            "post-1",
		Title:         "Test Post",
		Content:       "This is a test post.",
		AuthorID:      "user-1",
		AllowComments: true,
		CreatedAt:     time.Now(),
	}
	createdPost, err := store.CreatePost(context.Background(), post)
	assert.NoError(t, err, "CreatePost should not return an error")
	assert.NotNil(t, createdPost, "Created post should not be nil")
	assert.NotEmpty(t, createdPost.ID, "Post ID should not be empty")
	assert.Equal(t, post.ID, createdPost.ID, "Post IDs should match")
	assert.Equal(t, post.Title, createdPost.Title, "Post titles should match")
	assert.Equal(t, post.Content, createdPost.Content, "Post content should match")
	assert.Equal(t, post.AuthorID, createdPost.AuthorID, "Post author IDs should match")
	assert.Equal(t, post.AllowComments, createdPost.AllowComments, "Post allow comments flag should match")
	assert.Equal(t, post.CreatedAt, createdPost.CreatedAt, "Post creation times should match")
}

func TestInMemoryAddComment(t *testing.T) {
	store := NewStorageInMemory()
	post := &models.Post{
		ID:            "post-1",
		Title:         "Test Post",
		Content:       "This is a test post.",
		AuthorID:      "user-1",
		AllowComments: true,
		CreatedAt:     time.Now(),
	}
	store.CreatePost(context.Background(), post)

	comment := &models.Comment{
		PostID:    "post-1",
		ParentID:  nil,
		AuthorID:  "user-2",
		Text:      "Test comment",
		CreatedAt: time.Now(),
	}
	addedComment, err := store.AddComment(context.Background(), comment)
	assert.NoError(t, err, "AddComment should not return an error")
	assert.NotNil(t, addedComment, "Added comment should not be nil")
	assert.NotEmpty(t, addedComment.ID, "Comment ID should not be empty")
	assert.Equal(t, comment.PostID, addedComment.PostID, "Comment post IDs should match")
	assert.Equal(t, comment.ParentID, addedComment.ParentID, "Comment parent IDs should match")
	assert.Equal(t, comment.AuthorID, addedComment.AuthorID, "Comment author IDs should match")
	assert.Equal(t, comment.Text, addedComment.Text, "Comment texts should match")
	assert.Equal(t, comment.CreatedAt, addedComment.CreatedAt, "Comment creation times should match")
}

func TestInMemoryGetPosts(t *testing.T) {
	store := NewStorageInMemory()
	for i := 1; i < 11; i++ {
		post := &models.Post{
			Title:         fmt.Sprintf("Post %d", i),
			Content:       fmt.Sprintf("Test text %d", i),
			AuthorID:      "user-1",
			AllowComments: true,
			CreatedAt:     time.Now().Add(time.Duration(-i) * time.Minute),
		}
		store.CreatePost(context.Background(), post)
		time.Sleep(50 * time.Millisecond)
	}

	postsList, err := store.GetPosts(context.Background())
	assert.NoError(t, err, "GetPosts should not return an error")
	assert.Len(t, postsList, 10, "Should return exactly 10 posts")

	for i := 0; i < len(postsList)-1; i++ {
		assert.True(t, postsList[i].CreatedAt.After(postsList[i+1].CreatedAt), "Posts should be sorted by created_at in descending order")
	}
}

func TestInMemoryGetLatestComment(t *testing.T) {
	store := NewStorageInMemory()
	post := &models.Post{
		ID:            "post-1",
		Title:         "Test Post",
		Content:       "This is a test post.",
		AuthorID:      "user-1",
		AllowComments: true,
		CreatedAt:     time.Now(),
	}
	store.CreatePost(context.Background(), post)

	var now time.Time
	for i := 1; i < 5; i++ {
		comment := &models.Comment{
			ID:        fmt.Sprintf("com-%d", i),
			PostID:    "post-1",
			ParentID:  nil,
			AuthorID:  "user-2",
			Text:      fmt.Sprintf("Test tcomment %d", i),
			CreatedAt: time.Now(),
		}
		now = time.Now()
		store.AddComment(context.Background(), comment)
		time.Sleep(50 * time.Millisecond)
	}

	latestComment, err := store.GetLatestComment(context.Background(), post.ID)
	assert.NoError(t, err, "GetLatestComment should not return an error")
	assert.NotNil(t, latestComment, "Latest comment should not be nil")
	assert.WithinDuration(t, now, latestComment.CreatedAt, 1*time.Millisecond, "Latest comment's creation time should be close to the current time")

	childComment := &models.Comment{
		PostID:    "post-1",
		ParentID:  &latestComment.ID,
		AuthorID:  "user-3",
		Text:      "Test tcomment child",
		CreatedAt: time.Now(),
	}
	store.AddComment(context.Background(), childComment)

	latestComment, err = store.GetLatestComment(context.Background(), post.ID)
	assert.NoError(t, err, "GetLatestComment should not return an error")
	assert.Nil(t, latestComment.ParentID, "Latest comment should be a parent comment")
}

func TestInMemoryGetPostByID(t *testing.T) {
	store := NewStorageInMemory()
	post := &models.Post{
		ID:            "post-1",
		Title:         "Test Post",
		Content:       "This is a test post.",
		AuthorID:      "user-1",
		AllowComments: true,
		CreatedAt:     time.Now(),
	}
	store.CreatePost(context.Background(), post)

	receivedPost, err := store.GetPostByID(context.Background(), post.ID)
	assert.NoError(t, err, "GetPostByID should not return an error")
	assert.NotNil(t, receivedPost, "Received post should not be nil")
	assert.Equal(t, post.ID, receivedPost.ID, "Post IDs should match")
	assert.Equal(t, post.Title, receivedPost.Title, "Post titles should match")
	assert.Equal(t, post.Content, receivedPost.Content, "Post content should match")
	assert.Equal(t, post.AuthorID, receivedPost.AuthorID, "Post author IDs should match")
	assert.Equal(t, post.AllowComments, receivedPost.AllowComments, "Post allow comments flag should match")
	assert.Equal(t, post.CreatedAt, receivedPost.CreatedAt, "Post creation times should match")
}
