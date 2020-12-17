package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/eviltomorrow/aphrodite-base/zlog"
	//
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

//
var (
	MySQLDSN       string
	MySQLMinOpen   int = 5
	MySQLMaxOpen   int = 10
	MySQL          *sql.DB
	DefaultTimeout = 10 * time.Second
)

// BuildMySQL build mysql
func BuildMySQL() {
	pool, err := buildMySQL(MySQLDSN)
	if err != nil {
		zlog.Fatal("Build mysql connection failure", zap.Error(err))
	}
	MySQL = pool
}

// CloseMySQL close mysql
func CloseMySQL() error {
	zlog.Info("Close mysql connection", zap.String("dsn", MySQLDSN))

	if MySQL == nil {
		return nil
	}

	err := MySQL.Close()
	if err != nil {
		zlog.Error("Close mysql connection failure", zap.Error(err))
	}
	return err
}

func buildMySQL(dsn string) (*sql.DB, error) {
	pool, err := sql.Open("mysql", MySQLDSN)
	if err != nil {
		return nil, err
	}
	pool.SetConnMaxLifetime(time.Minute * 3)
	pool.SetMaxOpenConns(MySQLMaxOpen)
	pool.SetMaxIdleConns(MySQLMinOpen)

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	if err = pool.PingContext(ctx); err != nil {
		return nil, err
	}
	return pool, nil
}

// ExecMySQL exec mysql
type ExecMySQL interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}
