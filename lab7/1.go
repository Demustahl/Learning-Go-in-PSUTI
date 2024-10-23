package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func main() {
	// Определяем порт для прослушивания
	port := ":8080"

	// Запускаем прослушивание указанного порта
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
		return
	}
	defer listener.Close()
	fmt.Println("TCP-сервер запущен на порту", port)

	for {
		// Принимаем входящее соединение
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Ошибка при принятии соединения:", err)
			continue
		}
		fmt.Println("Новое соединение принято")

		// Обрабатываем соединение
		handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Обрабатываем несколько сообщений от клиента
	for {
		// Читаем сообщение от клиента
		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("Клиент закрыл соединение")
				return
			}
			fmt.Println("Ошибка при чтении сообщения:", err)
			return
		}

		// Выводим сообщение на экран
		fmt.Println("Сообщение от клиента:", message)

		// Отправляем ответ клиенту
		response := "Сообщение получено\n"
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Ошибка при отправке ответа:", err)
			return
		}
	}
}
