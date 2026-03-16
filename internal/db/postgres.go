package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB() error {
	fmt.Println("Connecting to the database...")
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	// dsn is my connection string in here which I'm passing into the pgxpool.New function to create a new connection pool to my PostgreSQL database. The context.Background() is used to provide a context for the connection, which can be useful for managing timeouts and cancellations.
	if err != nil {
		return err
	}
	DB = pool
	return pool.Ping(context.Background())
	//means that after creating the connection pool, I'm using the Ping method to check if the connection to the database is successful. If the ping is successful, it will return nil, indicating that the connection is established. If there is an error during the ping, it will return that error.
}
