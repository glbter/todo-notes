package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"
	"time"
)

func Connect(ctx context.Context) (*pgx.Conn, error) {
	for i := 0; i < 2; i++ {
		conn, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
		if err == nil {
			return conn, nil
		}

		select {
		case <- time.After(time.Second):
			continue
		case <- ctx.Done():
			return nil, fmt.Errorf("context finished")
		}
	}

	return nil, fmt.Errorf("could not connect to database")
}

func Ping(ctx context.Context, conn *pgx.Conn) error {
	for i := 0; i < 2; i++ {
		if err := conn.Ping(ctx); err == nil {
			return nil
		}

		select {
		case <- time.After(time.Second):
			continue
		case <- ctx.Done():
			return fmt.Errorf("context finished")
		}
	}

	return fmt.Errorf("could not ping the database")
}
