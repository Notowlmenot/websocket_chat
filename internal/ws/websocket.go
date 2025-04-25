package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn *websocket.Conn
}

// Глобальная переменная для хранения списка клиентов
var clients = make(map[*Client]bool)
var mutex sync.Mutex // Mutex для защиты доступа к списку клиентов

func WebsocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}

	client := &Client{conn: conn}
	mutex.Lock()
	clients[client] = true
	mutex.Unlock()

	log.Println("Client connected")
	defer func() {
		mutex.Lock()
		delete(clients, client)
		mutex.Unlock()
		conn.Close()
		log.Println("Client disconnected")
	}()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		log.Printf("Received message: %s\n", p)

		mutex.Lock()
		clientsCopy := make([]*Client, 0, len(clients))
		for client := range clients {
			clientsCopy = append(clientsCopy, client)
		}
		mutex.Unlock()

		// Итерируем копию списка клиентов
		for _, client := range clientsCopy {
			if client.conn != conn {
				err := client.conn.WriteMessage(messageType, p)
				if err != nil {
					log.Println("Error writing message:", err)
					mutex.Lock()
					delete(clients, client) // Удаляем из оригинального списка
					mutex.Unlock()
					client.conn.Close()
				}
			}
		}
	}
}
