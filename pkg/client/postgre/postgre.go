package postgre

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type conf struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

var (
	us1 = conf{
		Username: "postgres",
		Password: "postgres",
		Host:     "localhost",
		Port:     "5432",
		Database: "",
	}
)

type Client interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

func NewClient(ctx context.Context, atempt int) (pool *pgxpool.Pool, err error) { //conf &config.StorageConfig
	con := fmt.Sprintf("postgresql://%s:%s@%s:%s", us1.Username, us1.Password, us1.Host, us1.Port)
	AttempingForConn(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.Connect(ctx, con)
		if err != nil {
			return err
		}
		return nil
	}, atempt, 5*time.Second)

	if err != nil {
		fmt.Println(err)
		log.Fatal("Err of AttempingForConn")
	}
	return
}
