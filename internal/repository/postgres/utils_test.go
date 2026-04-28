package repository_test

import (
	repository "sugdio/internal/repository/postgres"
	"testing"
)

func TestPatchBuilder(t *testing.T) {
	t.Run("success update with where", func(t *testing.T) {
		builder := repository.NewPatchBuilder()

		builder.Head("users")
		builder.Add("name", "Alice")
		builder.Add("age", 30)
		builder.Where("id", 1)

		expectedSQL := "UPDATE users SET name = $1, age = $2 WHERE id = $3"
		if got := builder.String(); got != expectedSQL {
			t.Errorf("SQL: got %q, want %q", got, expectedSQL)
		}

		expectedArgs := []any{"Alice", 30, 1}
		args := builder.Args()
		if len(args) != len(expectedArgs) {
			t.Fatalf("Args length: got %d, want %d", len(args), len(expectedArgs))
		}

		for i, v := range args {
			if v != expectedArgs[i] {
				t.Errorf("Arg[%d]: got %v, want %v", i, v, expectedArgs[i])
			}
		}
	})
}
