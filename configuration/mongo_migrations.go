package configuration

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type MongoMigrations struct {
	m *migrate.Migrate
}

func NewMongoMigrations() (*MongoMigrations, error) {
	m, err := migrate.New(
		"file://mongo_migrations",
		"mongodb://localhost:27017/gott")
	if err != nil {
		return nil, err
	}

	return &MongoMigrations{
		m: m,
	}, nil
}

func (migrations *MongoMigrations) MigrateAll() error {
	return migrations.m.Up()
}

func (migrations *MongoMigrations) DownAll() error {
	return migrations.m.Down()
}
