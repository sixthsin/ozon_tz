package storage

import (
	"context"
	"database/sql"
	"fmt"
	"ozontz/app/models"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpassword",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err, "Failed to start container")

	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err, "Failed to get mapped port")

	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		"testuser", "testpassword", "localhost", port.Int(), "testdb")

	time.Sleep(5 * time.Second)

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err, "Failed to open database connection")

	migrationDir := getMigrationDir()
	t.Logf("Using migration directory: %s", migrationDir)

	store := NewStoragePostgres(db)
	err = store.ApplyMigrations(migrationDir)
	require.NoError(t, err, "Failed to apply migrations")

	teardown := func() {
		db.Close()
		container.Terminate(ctx)
	}

	return db, teardown
}

func TestCreatePost(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()

	store := NewStoragePostgres(db)

	post := &models.Post{
		Title:         "Test Post",
		Content:       "Test text",
		AuthorID:      "user-1",
		AllowComments: true,
	}

	createdPost, err := store.CreatePost(context.Background(), post)
	require.NoError(t, err, "CreatePost failed")
	require.NotEmpty(t, createdPost.ID, "Created post ID is empty")
	require.Equal(t, post.Title, createdPost.Title, "Post title mismatch")
	require.Equal(t, post.AuthorID, createdPost.AuthorID, "Post author ID mismatch")
}

func TestGetPostByID(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()

	store := NewStoragePostgres(db)

	post := &models.Post{
		Title:         "Test Post",
		Content:       "Test text",
		AuthorID:      "user-1",
		AllowComments: true,
	}
	createdPost, err := store.CreatePost(context.Background(), post)
	require.NoError(t, err, "CreatePost failed")

	fetchedPost, err := store.GetPostByID(context.Background(), createdPost.ID)
	require.NoError(t, err, "GetPostByID failed")
	assert.Equal(t, createdPost.ID, fetchedPost.ID, "Post ID mismatch")
	assert.Equal(t, createdPost.Title, fetchedPost.Title, "Post title mismatch")
}

func TestAddComment(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()

	store := NewStoragePostgres(db)

	post := &models.Post{
		Title:         "Test Post",
		Content:       "Test text",
		AuthorID:      "user-1",
		AllowComments: true,
	}
	createdPost, err := store.CreatePost(context.Background(), post)
	require.NoError(t, err, "CreatePost failed")

	comment := &models.Comment{
		PostID:   createdPost.ID,
		AuthorID: "user-2",
		Text:     "Test comment",
	}

	createdComment, err := store.AddComment(context.Background(), comment)
	require.NoError(t, err, "AddComment failed")
	assert.NotEmpty(t, createdComment.ID, "Created comment ID is empty")
	assert.Equal(t, comment.Text, createdComment.Text, "Comment text mismatch")
}

func TestGetLatestComment(t *testing.T) {
	db, teardown := setupTestDB(t)
	defer teardown()

	store := NewStoragePostgres(db)

	post := &models.Post{
		Title:         "Test Post",
		Content:       "Test text",
		AuthorID:      "user-1",
		AllowComments: true,
	}
	createdPost, err := store.CreatePost(context.Background(), post)
	require.NoError(t, err, "CreatePost failed")

	var now time.Time
	for i := 1; i < 5; i++ {
		comment := &models.Comment{
			ID:        fmt.Sprintf("com-%d", i),
			PostID:    createdPost.ID,
			ParentID:  nil,
			AuthorID:  "user-2",
			Text:      fmt.Sprintf("Test comment %d", i),
			CreatedAt: time.Now().UTC(),
		}
		now = time.Now().UTC()
		store.AddComment(context.Background(), comment)
		time.Sleep(50 * time.Millisecond)
	}

	latestComment, err := store.GetLatestComment(context.Background(), createdPost.ID)
	require.NoError(t, err, "GetLatestComment failed")

	const tolerance = 1 * time.Millisecond
	if diff := latestComment.CreatedAt.Sub(now); diff > tolerance || diff < -tolerance {
		t.Fatalf("Comment not latest time: got %v, want %v (diff: %v)", latestComment.CreatedAt, now, diff)
	}

	childComment := &models.Comment{
		PostID:    createdPost.ID,
		ParentID:  &latestComment.ID,
		AuthorID:  "user-3",
		Text:      "Test comment child",
		CreatedAt: time.Now().UTC(),
	}
	store.AddComment(context.Background(), childComment)

	latestComment, err = store.GetLatestComment(context.Background(), createdPost.ID)
	require.NoError(t, err, "GetLatestComment failed")
	assert.Nil(t, latestComment.ParentID, "Comment not parent")
}
