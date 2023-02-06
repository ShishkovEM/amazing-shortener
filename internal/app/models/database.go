package models

import (
	"context"
	"errors"

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
		"short_url text COLLATE pg_catalog.\"default\" NOT NULL,"+
		"original_url text COLLATE pg_catalog.\"default\" NOT NULL,"+
		"user_id bigint NOT NULL,"+
		"created_at timestamp with time zone,"+
		"CONSTRAINT urls_pkey PRIMARY KEY (short_url)"+
		")")

	if err != nil {
		panic(err)
	}
}
