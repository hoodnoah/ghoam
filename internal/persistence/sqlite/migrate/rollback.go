package migrate

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	msqlite "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func Rollback(db *sql.DB, steps int) error {
	drv, err := msqlite.WithInstance(db, &msqlite.Config{})
	if err != nil {
		return fmt.Errorf("failed to create sqlite migration driver: %w", err)
	}

	src, err := iofs.New(migrations, "sql")
	if err != nil {
		return fmt.Errorf("failed to create source driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "sqlite3", drv)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Execute the rollback steps.
	// For example, if steps == -1, this will undo the most recent migration.
	if err := m.Steps(steps); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration rollback error: %w", err)
	}

	return nil
}
