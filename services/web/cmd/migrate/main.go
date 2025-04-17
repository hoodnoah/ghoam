package migrate

import (
	// internal
	"github.com/hoodnoah/ghoam/internal/persistence/sqlite"
)

// Execute bootstraps the database, runs migrations, and inserts the basic account groups
// Returns repositories which abstract data retrieval, saving, etc.
func Execute(dbPath string) (*sqlite.Repositories, error) {
	// Open the SQLite database
	repos, err := sqlite.New(dbPath)
	if err != nil {
		return nil, err
	}

	return repos, nil
}
