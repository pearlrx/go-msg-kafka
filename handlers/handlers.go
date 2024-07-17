package handlers

import (
	"encoding/json"
	"log"
	"microservice-on-go/kafka"
	"microservice-on-go/messages"
	"net/http"
)

type MessageHandler struct {
	store    *messages.MessageStore
	producer *kafka.KafkaProducer
}

func NewMessageHandler(store *messages.MessageStore, producer *kafka.KafkaProducer) *MessageHandler {
	return &MessageHandler{store: store, producer: producer}
}

func (handler *MessageHandler) HandlePostMessage(w http.ResponseWriter, r *http.Request) {
	var message struct {
		Text string `json:"text"`
	}
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = handler.store.AddMessage(message.Text)
	if err != nil {
		log.Printf("Failed to store message: %v", err)
		http.Error(w, "Failed to store message", http.StatusInternalServerError)
		return
	}

	err = handler.producer.SendMessage(message.Text)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
		http.Error(w, "Failed to send message to Kafka", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Message received and stored"})
}

func (handler *MessageHandler) HandleGetStats(w http.ResponseWriter, r *http.Request) {
	count, err := handler.store.GetMessageCount()
	if err != nil {
		log.Printf("Failed to retrieve stats: %v", err)
		http.Error(w, "Failed to retrieve stats", http.StatusInternalServerError)
		return
	}

	stats := struct {
		MessageCount int `json:"message_count"`
	}{
		MessageCount: count,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
