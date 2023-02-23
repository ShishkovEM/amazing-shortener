package models

import (
	"context"
	"database/sql"
	"errors"
	"net/url"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
)

type DB struct {
	dsn  string
	conn *pgx.Conn
}

func NewDB(dsn string) *DB {
	return &DB{dsn: dsn}
}

func (d *DB) GetConnection() *pgx.Conn {
	return d.conn
}

func (d *DB) GetConn(ctx context.Context) (*pgx.Conn, error) {
	if d.dsn == "" {
		return nil, errors.New("empty database dsn")
	}

	conn, err := pgx.Connect(ctx, d.dsn)

	if err != nil {
		return nil, err
	}

	d.conn = conn

	return d.conn, nil
}

func (d DB) Close() {
	if d.conn == nil {
		return
	}

	err := d.conn.Close(context.Background())
	if err != nil {
		panic(err)
	}
}

func (d *DB) Migrate() error {
	db, err := sql.Open("postgres", d.dsn)
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	databasePath, err := url.Parse(d.dsn)
	if err != nil {
		panic(err)
	}
	databaseName := databasePath.Path[1:]

	m, mErr := migrate.NewWithDatabaseInstance(
		"file://./schema",
		databaseName, driver,
	)
	if mErr != nil {
		panic(mErr)
	}

	errOnMigrate := m.Up()

	if errOnMigrate != nil && !errors.Is(errOnMigrate, migrate.ErrNoChange) {
		panic(errOnMigrate)
	}

	return nil
}
