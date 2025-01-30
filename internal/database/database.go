package database

import "audit/internal/models"

type Database interface {
	GetLatestAuditLog() (*models.AuditLog, error)
}
