package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"lab6server/internal/domen/models"

	"github.com/gorilla/websocket"
)

type ChatUsecase interface {
	AddMessage(login string, password string, text string) (*models.Message, error)
	GetChatHistory() ([]*models.Message, error)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ChatHandler struct {
	chatUsecase ChatUsecase
	clients     map[*websocket.Conn]struct{}
	mutex       sync.Mutex
}

func NewChatHandler(chatUsecase ChatUsecase) *ChatHandler {
	return &ChatHandler{
		chatUsecase: chatUsecase,
		clients:     make(map[*websocket.Conn]struct{}),
	}
}

func (h *ChatHandler) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	h.mutex.Lock()
	h.clients[conn] = struct{}{}
	h.mutex.Unlock()

	history, err := h.chatUsecase.GetChatHistory()
	if err != nil {
		log.Println("Ошибка при получении истории чата:", err)
		return
	}
	historyJSON, err := json.Marshal(history)
	if err != nil {
		log.Println("Ошибка при сериализации истории сообщений:", err)
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, historyJSON); err != nil {
		log.Println("Ошибка при отправке истории сообщений:", err)
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			h.mutex.Lock()
			delete(h.clients, conn)
			h.mutex.Unlock()
			break
		}

		parts := bytes.Split(msg, []byte(":"))
		if len(parts) != 3 {
			log.Println("Invalid message format")
			continue
		}
		login := string(parts[0])
		password := string(parts[1])
		text := string(parts[2])

		message, err := h.chatUsecase.AddMessage(login, password, text)
		if err != nil {
			errMsg := fmt.Sprintf("Error adding message: %v", err)
			if writeErr := conn.WriteMessage(websocket.TextMessage, []byte(errMsg)); writeErr != nil {
				log.Println("Ошибка при отправке сообщения об ошибке:", writeErr)
			}
			continue
		}

		messageJSON, err := json.Marshal(message)
		if err != nil {
			log.Println("Ошибка при сериализации сообщения:", err)
			continue
		}

		h.broadcastMessage(messageJSON)
	}
}

func (h *ChatHandler) broadcastMessage(message []byte) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for client := range h.clients {
		if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Ошибка при отправке сообщения: %v", err)
			client.Close()
			delete(h.clients, client)
		}
	}
}
