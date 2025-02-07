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
		"comments":      &graphql.Field{Type: graphql.NewList(commentType)},
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
			Resolve: resolvePosts,
		},
		"post": &graphql.Field{
			Type: postType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: resolvePost,
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
	},
})
