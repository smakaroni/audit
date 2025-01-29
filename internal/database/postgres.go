package database

import (
	"audit/internal/models"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresDB struct {
	pool *pgxpool.Pool
}

func NewPostgresDB(connString string) (*PostgresDB, error) {
	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	return &PostgresDB{pool: pool}, nil
}

func (db *PostgresDB) CreateTableIfNotExists() error {
	_, err := db.pool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS audit_logs (
            id SERIAL PRIMARY KEY,
            proto_data BYTEA
        )
    `)
	return err
}

func (db *PostgresDB) InsertAuditLog(log *models.AuditLog) error {
	protoData, err := models.ToProto(log)
	if err != nil {
		return err
	}
	_, err = db.pool.Exec(context.Background(),
		"INSERT INTO audit_logs (proto_data) VALUES ($1)",
		protoData)
	return err
}

func (db *PostgresDB) Close() {
	db.pool.Close()
}
