package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
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

	// Читаем сообщение от пользователя
	fmt.Print("Введите сообщение для отправки: ")
	reader := bufio.NewReader(os.Stdin)
	message, _ := reader.ReadString('\n')

	// Отправляем сообщение на сервер
	_, err = conn.Write([]byte(message))
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
