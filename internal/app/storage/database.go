package storage

import (
	"context"

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
