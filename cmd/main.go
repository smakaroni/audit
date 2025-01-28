package main

import (
	"context"
	"log"

	"audit/internal/database"
	"audit/internal/kafka"
	"audit/internal/models"
)

func main() {
	// Set up Kafka consumer
	consumer, err := kafka.NewConsumer("localhost:9092", "audit_topic", "audit_group")
	if err != nil {
		log.Fatal("Failed to create Kafka consumer:", err)
	}
	defer consumer.Close()

	// Set up PostgreSQL connection
	db, err := database.NewPostgresDB("postgres://username:password@localhost:5432/audit_db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Create table if not exists
	if err := db.CreateTableIfNotExists(); err != nil {
		log.Fatal("Failed to create table:", err)
	}

	for {
		message, err := consumer.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error reading message:", err)
			continue
		}

		event, err := models.UnmarshalEvent(message.Value)
		if err != nil {
			log.Println("Error unmarshaling message:", err)
			continue
		}

		// Prepare data for insertion
		auditLog := models.NewAuditLog(event)

		// Insert into database
		if err := db.InsertAuditLog(auditLog); err != nil {
			log.Println("Error inserting into database:", err)
			continue
		}

		log.Printf("Saved anonymized event: %+v\n", event)
	}
}
