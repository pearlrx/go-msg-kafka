package messages

import (
	"database/sql"
	"log"
)

type MessageStore struct {
	db *sql.DB
}

func NewMessageStore(db *sql.DB) *MessageStore {
	return &MessageStore{db: db}
}

func (store *MessageStore) AddMessage(content string) error {
	_, err := store.db.Exec("INSERT INTO messages (content) VALUES ($1)", content)
	if err != nil {
		log.Printf("Error inserting message into database: %v", err)
	}
	return err
}

func (store *MessageStore) GetMessageCount() (int, error) {
	var count int
	err := store.db.QueryRow("SELECT COUNT(*) FROM messages").Scan(&count)
	if err != nil {
		log.Printf("Error retrieving message count from database: %v", err)
	}
	return count, err
}
