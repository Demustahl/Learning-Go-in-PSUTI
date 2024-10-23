package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gorilla/websocket"
)

func main() {
	// Подключаемся к серверу
	fmt.Print("Введите ваше имя: ")
	reader := bufio.NewReader(os.Stdin)
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	if username == "" {
		fmt.Println("Имя не может быть пустым.")
		return
	}

	// Устанавливаем соединение с веб-сокет сервером
	ws, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatalf("Не удалось подключиться к серверу: %v", err)
	}
	defer ws.Close()

	// Канал для обработки сигналов завершения
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Канал для отправки сообщений
	done := make(chan struct{})

	// Горутина для чтения сообщений от сервера
	go func() {
		defer close(done)
		for {
			var msg map[string]string
			err := ws.ReadJSON(&msg)
			if err != nil {
				log.Println("Ошибка при чтении сообщения:", err)
				return
			}
			fmt.Printf("%s: %s\n", msg["username"], msg["message"])
		}
	}()

	fmt.Println("Теперь вы можете отправлять сообщения. Введите 'exit' для выхода.")

	// Цикл для отправки сообщений
	for {
		fmt.Print("> ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		if message == "exit" {
			fmt.Println("Выход из чата...")
			break
		}

		// Отправляем сообщение на сервер
		err := ws.WriteJSON(map[string]string{
			"username": username,
			"message":  message,
		})
		if err != nil {
			log.Println("Ошибка при отправке сообщения:", err)
			break
		}
	}

	// Ожидаем завершения либо по сигналу, либо после выхода из цикла
	select {
	case <-done:
	case <-interrupt:
		fmt.Println("Получен сигнал завершения, закрытие соединения...")
		ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		ws.Close()
	}
}
