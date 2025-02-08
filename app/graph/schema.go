package graph

import (
	"github.com/graphql-go/graphql"
)

var postType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Post",
	Fields: graphql.Fields{
		"id":            &graphql.Field{Type: graphql.String},
		"title":         &graphql.Field{Type: graphql.String},
		"content":       &graphql.Field{Type: graphql.String},
		"authorId":      &graphql.Field{Type: graphql.String},
		"allowComments": &graphql.Field{Type: graphql.Boolean},
		"createdAt":     &graphql.Field{Type: graphql.String},
		"lastComment": &graphql.Field{
			Type:    commentType,
			Resolve: resolveGetLastComment,
		},
	},
})

var commentType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Comment",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.String},
		"postId":    &graphql.Field{Type: graphql.String},
		"parentId":  &graphql.Field{Type: graphql.String},
		"authorId":  &graphql.Field{Type: graphql.String},
		"text":      &graphql.Field{Type: graphql.String},
		"createdAt": &graphql.Field{Type: graphql.String},
	},
})

var QueryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"posts": &graphql.Field{
			Type:    graphql.NewList(postType),
			Resolve: resolveGetPostsList,
		},
		"post": &graphql.Field{
			Type: postType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: resolveGetPost,
		},
	},
})

var MutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"createPost": &graphql.Field{
			Type: postType,
			Args: graphql.FieldConfigArgument{
				"title":         &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"content":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"authorId":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"allowComments": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Boolean)},
			},
			Resolve: resolveCreatePost,
		},
		"addComment": &graphql.Field{
			Type: commentType,
			Args: graphql.FieldConfigArgument{
				"postId":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"parentId": &graphql.ArgumentConfig{Type: graphql.String},
				"authorId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"text":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: resolveAddComment,
		},
	},
})
