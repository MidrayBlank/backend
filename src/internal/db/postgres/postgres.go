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

type PostgresDBTransaction struct {
	tx *gorm.DB
}

func NewPostgresConnection(dsn string) PostgresDBConnection {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("Can not connect to postgres DB: %s", err.Error()))
	}

	return PostgresDBConnection{conn: db}
}

func (db PostgresDBConnection) Get() any {
	return db.conn
}

func (db PostgresDBConnection) BeginTx() abstract.IDBTransaction {
	return PostgresDBTransaction{tx: db.conn.Begin()}
}

func (db PostgresDBTransaction) Commit() error {
	return db.tx.Commit().Error
}

func (db PostgresDBTransaction) Rollback() {
	db.tx.Rollback()
}
