package accounting

import (
	"database/sql"
	"testing"
)

func TestAccountNewGroup(t *testing.T) {
	t.Run("fails with a blank group name", func(t *testing.T) {
		_, err := NewAccountGroup("", "Parent Group", sql.NullString{})
		if err == nil {
			t.Fatalf("expected an error, didn't receive one")
		}
	})

	t.Run("fails with a blank parent group name", func(t *testing.T) {
		_, err := NewAccountGroup("Group Name", "", sql.NullString{})
		if err == nil {
			t.Fatalf("expected an error, didn't receive one")
		}
	})

	t.Run("fails with a valid, but blank DisplayAfter name", func(t *testing.T) {
		_, err := NewAccountGroup("Group Name", "Parent Name", sql.NullString{String: "", Valid: true})
		if err == nil {
			t.Fatalf("expected an error, didn't receive one")
		}
	})
}
