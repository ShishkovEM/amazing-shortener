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

func (d *DB) CreateTables() {
	conn, err := d.GetConn(context.Background())

	defer d.Close()

	if err != nil {
		panic(err)
	}

	_, err = conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS urls ("+
		"short_uri text COLLATE pg_catalog.\"default\" NOT NULL,"+
		"original_url text COLLATE pg_catalog.\"default\" NOT NULL,"+
		"user_id bigint NOT NULL,"+
		"created_at timestamp with time zone,"+
		"CONSTRAINT urls_pkey PRIMARY KEY (original_url)"+
		")")

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

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/app/migrations",
		databaseName, driver,
	)
	if err != nil {
		panic(err)
	}

	errOnMigrate := m.Up()

	if errOnMigrate != nil && !errors.Is(errOnMigrate, migrate.ErrNoChange) {
		panic(errOnMigrate)
	}

	return nil
}
