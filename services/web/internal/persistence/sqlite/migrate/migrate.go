package migrate

import (
	"database/sql"
	"embed"

	"github.com/golang-migrate/migrate/v4"
	msqlite "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed sql/*.sql
var migrations embed.FS

func Up(db *sql.DB) error {
	drv, err := msqlite.WithInstance(db, &msqlite.Config{})
	if err != nil {
		return err
	}

	src, err := iofs.New(migrations, "sql")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", src, "sqlite3", drv)
	if err != nil {
		return err
	}

	// idempotent: if already at latest, this is a no-op
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
