package main

import (
	"lab6server/internal/domen/usecase"
	"lab6server/internal/handlers"
	"lab6server/internal/services"
	messagestorage "lab6server/internal/storages/message_storage"
	usersstorage "lab6server/internal/storages/users_storage"
	"log"
	"net/http"
)

func main() {
	userStorage := usersstorage.NewUserStorage()
	messageStore := messagestorage.NewMessageStorage()
	authService := services.NewAuthService(userStorage)
	chatUsecase := usecase.NewChatUsecase(authService, messageStore)

	chatHandler := handlers.NewChatHandler(chatUsecase)

	http.HandleFunc("/ws", chatHandler.HandleConnections)

	cleanupService := services.NewMessageCleanupService(messageStore)
	cleanupService.Start()
	defer cleanupService.Stop()

	log.Println("Starting server on :80")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
