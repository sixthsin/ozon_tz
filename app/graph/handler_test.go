package graph

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
)

var mockSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"hello": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "world", nil
				},
			},
		},
	}),
})

func TestGraphQLHandler(t *testing.T) {
	handler := GraphQLHandler(&mockSchema)

	t.Run("Successful request", func(t *testing.T) {
		requestBody := `{"query": "{ hello }"}`

		req := httptest.NewRequest(http.MethodPost, "/query", bytes.NewBufferString(requestBody))
		w := httptest.NewRecorder()

		handler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		data, ok := response["data"].(map[string]interface{})
		assert.True(t, ok)
		hello, ok := data["hello"].(string)
		assert.True(t, ok)
		assert.Equal(t, "world", hello)
	})

	t.Run("Invalid JSON payload", func(t *testing.T) {
		requestBody := `{"query": "invalid"`

		req := httptest.NewRequest(http.MethodPost, "/query", bytes.NewBufferString(requestBody))
		w := httptest.NewRecorder()

		handler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		assert.Contains(t, w.Body.String(), "Invalid JSON payload")
	})

	t.Run("GraphQL errors", func(t *testing.T) {
		requestBody := `{"query": "{ invalidField }"}`

		req := httptest.NewRequest(http.MethodPost, "/query", bytes.NewBufferString(requestBody))
		w := httptest.NewRecorder()

		handler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		assert.Contains(t, w.Body.String(), "GraphQL errors")
	})

	t.Run("Missing request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/query", nil)
		w := httptest.NewRecorder()

		handler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		assert.Contains(t, w.Body.String(), "Request body is empty")
	})
}
