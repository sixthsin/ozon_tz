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

func resolveGetPostsList(params graphql.ResolveParams) (interface{}, error) {
	posts, err := store.GetPosts(context.Background())
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func resolveGetPost(params graphql.ResolveParams) (interface{}, error) {
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
	title := params.Args["title"].(string)
	content := params.Args["content"].(string)
	authorId := params.Args["authorId"].(string)
	allowComments := params.Args["allowComments"].(bool)

	post := &models.Post{
		Title:         title,
		Content:       content,
		AuthorID:      authorId,
		AllowComments: allowComments,
		CreatedAt:     time.Now(),
	}
	return store.CreatePost(context.Background(), post)
}

func resolveAddComment(params graphql.ResolveParams) (interface{}, error) {
	rawPostID := params.Args["postId"]
	if rawPostID == nil {
		return nil, errors.New("postId is required")
	}
	postId, ok := rawPostID.(string)
	if !ok {
		return nil, errors.New("postId must be a string")
	}

	rawParentID := params.Args["parentId"]
	var parentID *string
	if rawParentID != nil {
		parentIdStr, ok := rawParentID.(string)
		if !ok {
			return nil, errors.New("parentId must be a string")
		}
		parentID = &parentIdStr
	}

	rawAuthorID := params.Args["authorId"]
	if rawAuthorID == nil {
		return nil, errors.New("authorId is required")
	}
	authorId, ok := rawAuthorID.(string)
	if !ok {
		return nil, errors.New("authorId must be a string")
	}

	rawText := params.Args["text"]
	if rawText == nil {
		return nil, errors.New("text is required")
	}
	text, ok := rawText.(string)
	if !ok {
		return nil, errors.New("text must be a string")
	}

	comment := &models.Comment{
		PostID:    postId,
		ParentID:  parentID,
		AuthorID:  authorId,
		Text:      text,
		CreatedAt: time.Now(),
	}

	return store.AddComment(context.Background(), comment)
}

func resolveGetLastComment(params graphql.ResolveParams) (interface{}, error) {
	post, ok := params.Source.(*models.Post)
	if !ok {
		return nil, errors.New("invalid source type")
	}

	if !post.AllowComments {
		return nil, nil
	}

	lastComment, err := store.GetLatestComment(context.Background(), post.ID)
	if err != nil {
		return nil, err
	}

	return lastComment, nil
}
