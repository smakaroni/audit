package database

import (
	"audit/internal/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestInsertAuditLog(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mariaDB := &MariaDB{db: db}

	auditLog := &models.AuditLog{
		EventType:      "test",
		AnonymizedUser: "user123",
		Timestamp:      timestamppb.New(time.Now()),
		Data:           "test data",
	}

	protoData, _ := models.ToProto(auditLog)

	mock.ExpectExec("INSERT INTO audit_logs").
		WithArgs(protoData).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = mariaDB.InsertAuditLog(auditLog)
	if err != nil {
		t.Errorf("error was not expected while inserting audit log: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
