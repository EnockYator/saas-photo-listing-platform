package postgres

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const migrationsPath = "file://internal/infrastructure/database/postgres/migrations"

type Migrator struct {
	m *migrate.Migrate
}

func NewMigrator(databaseURL string) (*Migrator, error) {
	m, err := migrate.New(
		migrationsPath,
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
