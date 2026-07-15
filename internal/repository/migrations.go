package repository

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dbURL string) error {
	m, err := migrate.New("file:///app/migrations", dbURL)
	if err != nil {
		return fmt.Errorf("Ошибка создания мигратора: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("Ошибка применения миграций: %w", err)
	}

	log.Println("Миграции успешно применены")
	return nil
}
