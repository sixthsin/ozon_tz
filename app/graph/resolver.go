package graph

import (
	"context"
	"errors"
	"ozontz/app/models"
	"ozontz/app/storage"
	"time"

	"github.com/graphql-go/graphql"
)

var store storage.Storage

func SetStore(s storage.Storage) {
	store = s
}

func resolvePosts(params graphql.ResolveParams) (interface{}, error) {
	posts, err := store.GetPosts(context.Background())
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func resolvePost(params graphql.ResolveParams) (interface{}, error) {
	id, ok := params.Args["id"].(string)
	if !ok {
		return nil, errors.New("invalid ID")
	}
	post, err := store.GetPostByID(context.Background(), id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func resolveCreatePost(params graphql.ResolveParams) (interface{}, error) {
	title, _ := params.Args["title"].(string)
	content, _ := params.Args["content"].(string)
	authorId, _ := params.Args["authorId"].(string)
	allowComments, _ := params.Args["allowComments"].(bool)

	post := &models.Post{
		Title:         title,
		Content:       content,
		AuthorID:      authorId,
		AllowComments: allowComments,
		CreatedAt:     time.Now(),
	}
	return store.CreatePost(context.Background(), post)
}
