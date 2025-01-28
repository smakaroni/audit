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
			event_type TEXT,
			anonymized_user TEXT,
			timestamp TIMESTAMP,
			anonymized_data JSONB
		)
	`)
	return err
}

func (db *PostgresDB) InsertAuditLog(log *models.AuditLog) error {
	_, err := db.pool.Exec(context.Background(),
		"INSERT INTO audit_logs (event_type, anonymized_user, timestamp, anonymized_data) VALUES ($1, $2, $3, $4)",
		log.EventType, log.AnonymizedUser, log.Timestamp, log.Data)
	return err
}

func (db *PostgresDB) Close() {
	db.pool.Close()
}
