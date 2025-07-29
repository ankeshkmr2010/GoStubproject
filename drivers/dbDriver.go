package drivers

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"strconv"
)

func InitDB(ctx context.Context, config *Configuration) (*pgxpool.Pool, func()) {
	db := getConnectionForDb(ctx, config)
	return db, db.Close
}

func getConnectionForDb(ctx context.Context, config *Configuration) *pgxpool.Pool {
	user := config.DbUser
	host := config.DbHost
	port, err := strconv.Atoi(config.DbPort)
	if err != nil {
		port = 5432
	}

	password := config.DbPass
	dbName := config.DbName

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		panic(err)
	}

	db, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		panic(err)
	}

	err = db.Ping(ctx)
	if err != nil {
		panic(err)
	}

	log.Println("Database connection established")
	return db
}
