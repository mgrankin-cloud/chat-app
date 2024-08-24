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
	// Пример: postgres://user:password@localhost:5432/dbname?sslmode=disable
	flag.StringVar(&dbURL, "db-url", "", "URL for database connection")
	// Путь до папки с миграциями.
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	// Таблица, в которой будет храниться информация о миграциях. Она нужна
	// для того, чтобы понимать, какие миграции уже применены, а какие нет.
	// Дефолтное значение - 'migrations'.
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.Parse() // Выполняем парсинг флагов

	// Валидация параметров
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