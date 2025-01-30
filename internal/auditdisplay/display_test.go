package auditdisplay

import (
	"audit/internal/models"
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"os"
	"testing"
)

// MockDatabase is a mock of the Database interface
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) GetLatestAuditLog() (*models.AuditLog, error) {
	args := m.Called()
	return args.Get(0).(*models.AuditLog), args.Error(1)
}

func TestDisplayChanges(t *testing.T) {
	tests := []struct {
		name           string
		currentLog     *models.AuditLog
		previousLog    *models.AuditLog
		dbError        error
		expectedOutput string
		expectedError  string
	}{
		{
			name: "No changes",
			currentLog: &models.AuditLog{
				EventType:      "test",
				AnonymizedUser: "user1",
				Timestamp:      timestamppb.Now(),
				Data:           `{"field1": "value1", "field2": 42}`,
			},
			previousLog: &models.AuditLog{
				EventType:      "test",
				AnonymizedUser: "user1",
				Timestamp:      timestamppb.Now(),
				Data:           `{"field1": "value1", "field2": 42}`,
			},
			dbError:        nil,
			expectedOutput: "No changes detected.\n",
			expectedError:  "",
		},
		{
			name: "Changes detected",
			currentLog: &models.AuditLog{
				EventType:      "test",
				AnonymizedUser: "user1",
				Timestamp:      timestamppb.Now(),
				Data:           `{"field1": "new_value", "field2": 43, "field3": "added"}`,
			},
			previousLog: &models.AuditLog{
				EventType:      "test",
				AnonymizedUser: "user1",
				Timestamp:      timestamppb.Now(),
				Data:           `{"field1": "old_value", "field2": 42, "field4": "removed"}`,
			},
			dbError:        nil,
			expectedOutput: "Changes detected:\nField: field1\n  Old value: old_value\n  New value: new_value\nField: field2\n  Old value: 42\n  New value: 43\nField: field3\n  Old value: <nil>\n  New value: added\nField: field4\n  Old value: removed\n  New value: <nil>\n",
			expectedError:  "",
		},
		{
			name: "Database error",
			currentLog: &models.AuditLog{
				EventType:      "test",
				AnonymizedUser: "user1",
				Timestamp:      timestamppb.Now(),
				Data:           `{"field1": "value1"}`,
			},
			previousLog:    nil,
			dbError:        errors.New("database error"),
			expectedOutput: "",
			expectedError:  "error getting latest audit log: database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDatabase)
			mockDB.On("GetLatestAuditLog").Return(tt.previousLog, tt.dbError)

			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := DisplayChanges(mockDB, tt.currentLog)

			// Restore stdout
			w.Close()
			os.Stdout = oldStdout

			// Read captured output
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedOutput, output)

			mockDB.AssertExpectations(t)
		})
	}
}

func TestCompareAuditLogs(t *testing.T) {
	tests := []struct {
		name        string
		current     *models.AuditLog
		previous    *models.AuditLog
		expected    map[string]interface{}
		expectError bool
	}{
		{
			name:        "No changes",
			current:     &models.AuditLog{Data: `{"field1": "value1", "field2": 42}`},
			previous:    &models.AuditLog{Data: `{"field1": "value1", "field2": 42}`},
			expected:    map[string]interface{}{},
			expectError: false,
		},
		{
			name:     "Value changed",
			current:  &models.AuditLog{Data: `{"field1": "new_value", "field2": 42}`},
			previous: &models.AuditLog{Data: `{"field1": "old_value", "field2": 42}`},
			expected: map[string]interface{}{
				"field1": map[string]interface{}{
					"old": "old_value",
					"new": "new_value",
				},
			},
			expectError: false,
		},
		{
			name:     "Field added",
			current:  &models.AuditLog{Data: `{"field1": "value1", "field2": 42, "field3": "new"}`},
			previous: &models.AuditLog{Data: `{"field1": "value1", "field2": 42}`},
			expected: map[string]interface{}{
				"field3": map[string]interface{}{
					"old": nil,
					"new": "new",
				},
			},
			expectError: false,
		},
		{
			name:     "Field removed",
			current:  &models.AuditLog{Data: `{"field1": "value1"}`},
			previous: &models.AuditLog{Data: `{"field1": "value1", "field2": 42}`},
			expected: map[string]interface{}{
				"field2": map[string]interface{}{
					"old": float64(42),
					"new": nil,
				},
			},
			expectError: false,
		},
		{
			name:        "Invalid JSON in current",
			current:     &models.AuditLog{Data: `{"field1": "value1", invalid}`},
			previous:    &models.AuditLog{Data: `{"field1": "value1"}`},
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Invalid JSON in previous",
			current:     &models.AuditLog{Data: `{"field1": "value1"}`},
			previous:    &models.AuditLog{Data: `{"field1": "value1", invalid}`},
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := models.CompareAuditLogs(tt.current, tt.previous)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
