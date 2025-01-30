package auditdisplay

import (
	"audit/internal/database"
	"audit/internal/models"
	"fmt"
)

func DisplayChanges(db database.Database, currentLog *models.AuditLog) error {
	previousLog, err := db.GetLatestAuditLog()
	if err != nil {
		return fmt.Errorf("error getting latest audit log: %v", err)
	}

	changes, err := models.CompareAuditLogs(currentLog, previousLog)
	if err != nil {
		return fmt.Errorf("error comparing audit logs: %v", err)
	}

	if len(changes) == 0 {
		fmt.Println("No changes detected.")
		return nil
	}

	fmt.Println("Changes detected:")
	for key, change := range changes {
		changeMap := change.(map[string]interface{})
		fmt.Printf("Field: %s\n", key)
		fmt.Printf("  Old value: %v\n", changeMap["old"])
		fmt.Printf("  New value: %v\n", changeMap["new"])
	}

	return nil
}
