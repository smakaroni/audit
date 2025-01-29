package database

import (
	"audit/internal/models"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type MariaDB struct {
	db *sql.DB
}

func NewMariaDB(dsn string) (*MariaDB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &MariaDB{db: db}, nil
}

func (m *MariaDB) CreateTableIfNotExists() error {
	_, err := m.db.Exec(`
        CREATE TABLE IF NOT EXISTS audit_logs (
            id INT AUTO_INCREMENT PRIMARY KEY,
            proto_data BLOB
        )
    `)
	return err
}

func (m *MariaDB) InsertAuditLog(log *models.AuditLog) error {
	protoData, err := models.ToProto(log)
	if err != nil {
		return err
	}
	_, err = m.db.Exec(
		"INSERT INTO audit_logs (proto_data) VALUES (?)",
		protoData)
	return err
}

func (m *MariaDB) Close() error {
	return m.db.Close()
}
