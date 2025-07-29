package drivers

import (
	"context"
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	migratePGX "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"log"
)

func RunMigrate(ctx context.Context, mainDbPool *pgxpool.Pool, migrationFilePath string) {
	connString := mainDbPool.Config().ConnString()
	dbStd, err := sql.Open("pgx", connString)
	if err != nil {
		log.Fatal("Error preparing *sql.DB: ", err)
	}
	defer func() {
		if err := dbStd.Close(); err != nil {
			log.Fatal(ctx, "Error closing database", "error", err)
		}
	}()

	if err := dbStd.Ping(); err != nil {
		log.Fatal("Error pinging database: ", err)
	}
	driver, err := migratePGX.WithInstance(dbStd, &migratePGX.Config{})
	if err != nil {
		log.Fatal("Error in creating migration driver: ", err)
	}

	if migrationFilePath == "" {
		migrationFilePath = "file://migrations"
	}

	m, err := migrate.NewWithDatabaseInstance(migrationFilePath, "postgres", driver)
	if err != nil {
		log.Fatal("Error in creating migration: ", err)
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("Error in migration up: ", err)
	}
}
