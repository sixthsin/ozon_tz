package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
)

func GraphQLHandler(schema *graphql.Schema) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Установка заголовков
		w.Header().Set("Content-Type", "application/json")

		// Чтение тела запроса
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Разбор JSON-запроса
		var params struct {
			Query         string                 `json:"query"`
			OperationName string                 `json:"operationName"`
			Variables     map[string]interface{} `json:"variables"`
		}
		if err := json.Unmarshal(body, &params); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		// Выполнение GraphQL-запроса
		result := graphql.Do(graphql.Params{
			Context:        context.Background(),
			Schema:         *schema,
			RequestString:  params.Query,
			OperationName:  params.OperationName,
			VariableValues: params.Variables,
		})

		// Проверка ошибок
		if len(result.Errors) > 0 {
			log.Printf("GraphQL errors: %v\n", result.Errors)
			http.Error(w, fmt.Sprintf("GraphQL errors: %v", result.Errors), http.StatusBadRequest)
			return
		}

		// Отправка результата
		response, err := json.Marshal(result)
		if err != nil {
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			return
		}
		w.Write(response)
	}
}
