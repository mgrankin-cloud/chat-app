package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	// TODO: впихнуть инициализацию и подключить бд

	// TODO: сгенерить файлы ./storage/sso.db

	var dbURL, migrationsPath, migrationsTable string

	// Получаем необходимые значения из флагов запуска

	// URL для подключения к базе данных PostgreSQL.
	flag.StringVar(&dbURL, "db-url", "postgres://root:pass@localhost:5432/instant_messenger_db?sslmode=disable", "URL for database connection")
	// Путь до папки с миграциями.
	flag.StringVar(&migrationsPath, "migrations-path", "./migrations", "path to migrations")
	// Таблица, в которой будет храниться информация о миграциях.
	// Дефолтное значение - 'migrations'.
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.Parse() 

	if dbURL == "" {
		panic("db-url is required")
	}
	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	// Создаем объект мигратора, передав креды нашей БД
	m, err := migrate.New(
		"file://"+migrationsPath,
		dbURL)
	if err != nil {
		panic(err)
	}

	// Выполняем миграции до последней версии
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}

		panic(err)
	}
}
