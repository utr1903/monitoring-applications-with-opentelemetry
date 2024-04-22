package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/utr1903/monitoring-applications-with-opentelemetry/apps/commons/pkg/data"
)

type Database struct {
	instance *bun.DB
}

func NewDatabase() *Database {
	return &Database{}
}

func (d *Database) Connect(username string, password string, address string, port string, database string) {

	// dsn := "postgres://root:root@postgres:5432/test_db?sslmode=disable"
	dsn := "postgres://" + username + ":" + password + "@" + address + ":" + port + "/" + database + "?sslmode=disable"

	fmt.Println("Connecting...")
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	if db == nil {
		fmt.Println("Not connected.")
	} else {
		fmt.Println("Connected.")
	}

	d.instance = db
}

func (d *Database) CreateTableIfNotExists(ctx context.Context) {
	_, err := d.instance.NewCreateTable().Model((*data.Task)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		fmt.Println("Failed to create table:", err)
		return
	}
}
