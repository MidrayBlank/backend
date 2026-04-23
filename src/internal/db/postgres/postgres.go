package postgres

import (
	"fmt"

	"backend/src/internal/db/abstract"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDBConnection struct {
	conn *gorm.DB
}

func NewPostgresConnection(dsn string) abstract.IDBConnection {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("Can not connect to postgres DB: %s", err.Error()))
	}

	return &PostgresDBConnection{conn: db}
}

func (db PostgresDBConnection) Get() any {
	return db.conn
}

func (c PostgresDBConnection) BeginTx() abstract.IDBConnection {
	return &PostgresDBConnection{conn: c.conn.Begin()}
}

func (c PostgresDBConnection) Commit() error {
	return c.conn.Commit().Error
}

func (c PostgresDBConnection) Rollback() {
	c.conn.Rollback()
}
