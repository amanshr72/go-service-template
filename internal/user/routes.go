package user

import (
	"encoding/json"
	"net/http"

	"github.com/graphql-go/graphql"
)

func RegisterRoutes(mux *http.ServeMux, svc Service) error {
	h := NewHandler(svc)

	// REST routes
	mux.HandleFunc("POST /api/v1/users", h.Create)
	mux.HandleFunc("GET /api/v1/users", h.GetAll)
	mux.HandleFunc("GET /api/v1/users/{id}", h.GetByID)
	mux.HandleFunc("PUT /api/v1/users/{id}", h.Update)
	mux.HandleFunc("DELETE /api/v1/users/{id}", h.Delete)
	mux.HandleFunc("GET /api/v1/users/active/{active}", h.GetByActive)
	mux.HandleFunc("GET /api/v1/users/count", h.GetCount)

	// GraphQL — single endpoint
	schema, err := NewSchema(svc)
	if err != nil {
		return err
	}
	mux.HandleFunc("POST /graphql", graphqlHandler(schema))

	return nil
}

func graphqlHandler(schema graphql.Schema) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Query     string                 `json:"query"`
			Variables map[string]interface{} `json:"variables"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		result := graphql.Do(graphql.Params{
			Schema:         schema,
			RequestString:  body.Query,
			VariableValues: body.Variables,
		})
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(result)
	}
}
