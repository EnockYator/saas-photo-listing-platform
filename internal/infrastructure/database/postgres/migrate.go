package postgres

import (
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type Migrator struct {
	m *migrate.Migrate
}

func NewMigrator(databaseURL string) (*Migrator, error) {
	driver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithSourceInstance(
		"iofs",
		driver,
		databaseURL,
	)
	if err != nil {
		return nil, err
	}

	return &Migrator{
		m: m,
	}, nil
}

func (m *Migrator) Up() error {
	err := m.m.Up()

	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	return err
}

func (m *Migrator) Down() error {
	err := m.m.Down()

	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	return err
}

func (m *Migrator) Force(version int) error {
	return m.m.Force(version)
}

func (m *Migrator) Version() (uint, bool, error) {
	return m.m.Version()
}
