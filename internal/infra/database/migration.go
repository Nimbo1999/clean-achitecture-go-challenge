package database

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migration struct {
	db         *sql.DB
	driverName string
}

func NewMigration(db *sql.DB, driverName string) *Migration {
	return &Migration{db, driverName}
}

func (m *Migration) getMigration() (*migrate.Migrate, error) {
	driver, err := mysql.WithInstance(m.db, &mysql.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	return migrate.NewWithDatabaseInstance("file://sql/migrations", "mysql", driver)
}

func (m *Migration) Up() error {
	migration, err := m.getMigration()
	if err != nil {
		return err
	}

	err = migration.Up()
	if err == migrate.ErrNoChange {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

func (m *Migration) Down() error {
	migration, err := m.getMigration()
	if err != nil {
		return err
	}

	if err = migration.Down(); err != nil {
		return err
	}
	return nil
}
