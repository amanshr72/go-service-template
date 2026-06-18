package user

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
)

// helper: builds schema with mock service, executes a query string
func runQuery(t *testing.T, query string) *graphql.Result {
	t.Helper()
	svc := NewService(NewMockRepository(), &MockNotifier{})
	// pre-seed one user so queries have data
	_, _ = svc.Create(CreateUserInput{Name: "Aman", Email: "aman@t.com"})

	schema, err := NewSchema(svc)
	assert.NoError(t, err)

	return graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
}

func TestResolver_QueryUsers(t *testing.T) {
	result := runQuery(t, `{ users { id name email } }`)
	assert.Empty(t, result.Errors)

	users := result.Data.(map[string]interface{})["users"].([]interface{})
	assert.Len(t, users, 1)
	assert.Equal(t, "Aman", users[0].(map[string]interface{})["name"])
}

func TestResolver_QueryUserByID(t *testing.T) {
	result := runQuery(t, `{ user(id: 1) { id name } }`)
	assert.Empty(t, result.Errors)

	u := result.Data.(map[string]interface{})["user"].(map[string]interface{})
	assert.Equal(t, "Aman", u["name"])
}

func TestResolver_QueryUserByID_NotFound(t *testing.T) {
	result := runQuery(t, `{ user(id: 999) { id name } }`)
	// graphql-go puts resolver errors in result.Errors, not panics
	assert.NotEmpty(t, result.Errors)
}

func TestResolver_QueryCount(t *testing.T) {
	result := runQuery(t, `{ userCount }`)
	assert.Empty(t, result.Errors)

	count := result.Data.(map[string]interface{})["userCount"].(int)
	assert.Equal(t, 1, count)
}

func TestResolver_MutationCreate(t *testing.T) {
	svc := NewService(NewMockRepository(), &MockNotifier{})
	schema, _ := NewSchema(svc)

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `mutation { createUser(name: "Bob", email: "bob@t.com") { id name } }`,
	})
	assert.Empty(t, result.Errors)

	u := result.Data.(map[string]interface{})["createUser"].(map[string]interface{})
	assert.Equal(t, "Bob", u["name"])
}

func TestResolver_MutationDelete(t *testing.T) {
	svc := NewService(NewMockRepository(), &MockNotifier{})
	created, _ := svc.Create(CreateUserInput{Name: "Del", Email: "del@t.com"})
	schema, _ := NewSchema(svc)

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `mutation { deleteUser(id: ` + string(rune('0'+created.ID)) + `) }`,
	})
	assert.Empty(t, result.Errors)
}
