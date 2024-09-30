package tests

import (
	"testing"

	"github.com/pavozayac/constraint-service/tests/shared"
	"github.com/stretchr/testify/assert"
)

// example test for setting up CI
func TestTodo(t *testing.T) {
	client := shared.SetupClient()

	t.Run("Should get all todos", func(t *testing.T) {
		var response struct {
			Todos []struct {
				ID   string
				Text string
				Done bool
			}
		}

		client.MustPost(`
			query {
				todos {
					id text done
				}
			}
		`, &response)

		assert.Equal(t, response.Todos[0].ID, "1")
		assert.Equal(t, response.Todos[0].Text, "Some text")
		assert.Equal(t, response.Todos[0].Done, true)
	})

}
