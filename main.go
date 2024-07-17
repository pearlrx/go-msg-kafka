package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"

	"github.com/gorilla/mux"

	"microservice-on-go/handlers"
	"microservice-on-go/kafka"
	"microservice-on-go/messages"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		content TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	messageStore := messages.NewMessageStore(db)

	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		log.Fatalf("KAFKA_BROKERS environment variable is not set")
	}
	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		log.Fatalf("KAFKA_TOPIC environment variable is not set")
	}
	producer, err := kafka.NewKafkaProducer([]string{brokers}, topic)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	messageHandler := handlers.NewMessageHandler(messageStore, producer)

	r := mux.NewRouter()
	r.HandleFunc("/messages", messageHandler.HandlePostMessage).Methods("POST")
	r.HandleFunc("/stats", messageHandler.HandleGetStats).Methods("GET")

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
