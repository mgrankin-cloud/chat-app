package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Инициализация и подключение к базе данных
	// (предполагается, что база данных уже настроена)

	// Создание директории для хранения файлов, если она не существует
	storagePath := "./storage"
	if err := os.MkdirAll(storagePath, os.ModePerm); err != nil {
		panic(fmt.Errorf("failed to create storage directory: %w", err))
	}

	// Генерация файлов в директории ./storage/
	dbFilePath := filepath.Join(storagePath, "sso.db")
	if _, err := os.Create(dbFilePath); err != nil {
		panic(fmt.Errorf("failed to create db file: %w", err))
	}

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

/**package main

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
}**/
