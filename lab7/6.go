package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Объект для обновления соединений
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Разрешаем соединения с любого источника (для упрощения)
	},
}

// Хранение всех подключений
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var mutex = sync.Mutex{}

// Структура для сообщений
type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	http.HandleFunc("/ws", handleConnections)

	// Горутин для обработки входящих сообщений
	go handleMessages()

	// Запуск сервера на порту 8080
	fmt.Println("Сервер веб-сокетов запущен на порту 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Обновляем HTTP-соединение до веб-сокет-соединения
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Ошибка при обновлении соединения: %v", err)
		return
	}
	defer ws.Close()

	// Добавляем новое соединение в список клиентов
	mutex.Lock()
	clients[ws] = true
	mutex.Unlock()

	fmt.Println("Новое соединение установлено")

	// Чтение сообщений от клиента
	for {
		var msg Message
		// Читаем сообщение
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Ошибка при чтении сообщения: %v", err)
			mutex.Lock()
			delete(clients, ws)
			mutex.Unlock()
			fmt.Println("Соединение закрыто")
			break
		}

		// Отправляем сообщение в канал broadcast
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// Получаем сообщение из канала broadcast
		msg := <-broadcast

		// Отправляем сообщение всем подключенным клиентам
		mutex.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Ошибка при отправке сообщения: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}
