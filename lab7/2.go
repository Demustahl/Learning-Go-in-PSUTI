package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Определяем адрес сервера
	serverAddress := "localhost:8080"

	// Подключаемся к серверу
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("Ошибка при подключении к серверу:", err)
		return
	}
	defer conn.Close()

	// Создаем reader для чтения ввода пользователя
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Введите сообщение для отправки (или 'exit' для выхода):")

	for {
		// Читаем сообщение от пользователя
		fmt.Print("> ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		// Проверяем, если пользователь ввел "exit", то выходим из цикла
		if message == "exit" {
			fmt.Println("Закрытие соединения...")
			break
		}

		// Отправляем сообщение на сервер
		_, err = conn.Write([]byte(message + "\n"))
		if err != nil {
			fmt.Println("Ошибка при отправке сообщения:", err)
			return
		}

		// Ждем ответ от сервера
		responseReader := bufio.NewReader(conn)
		response, err := responseReader.ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка при чтении ответа от сервера:", err)
			return
		}

		// Выводим ответ сервера на экран
		fmt.Println("Ответ от сервера:", response)
	}
}
