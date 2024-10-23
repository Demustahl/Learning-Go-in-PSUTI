package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	port := ":8080"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
		return
	}
	fmt.Println("TCP-сервер запущен на порту", port)

	// Канал для отслеживания сигналов завершения
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Канал для завершения всех горутин
	ctx, cancel := context.WithCancel(context.Background())

	// Группа ожидания для отслеживания завершения всех горутин
	var wg sync.WaitGroup

	// Горутин для обработки новых соединений
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
					fmt.Println("Ошибка при принятии соединения:", err)
					continue
				}
			}

			fmt.Println("Новое соединение принято")

			// Увеличиваем счетчик горутин в группе ожидания
			wg.Add(1)
			go handleConnection(ctx, conn, &wg)
		}
	}()

	// Ожидаем сигнал завершения
	<-stopChan
	fmt.Println("\nПолучен сигнал завершения, завершаем сервер...")

	// Отменяем контекст, что остановит все активные горутины
	cancel()

	// Закрываем listener, чтобы больше не принимать новые соединения
	listener.Close()

	// Ждем завершения всех горутин
	wg.Wait()
	fmt.Println("Все соединения завершены, сервер остановлен.")
}

func handleConnection(ctx context.Context, conn net.Conn, wg *sync.WaitGroup) {
	defer conn.Close()
	defer wg.Done()

	reader := bufio.NewReader(conn)

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

		// Выводим сообщение от клиента
		fmt.Println("Сообщение от клиента:", message)

		// Отправляем ответ клиенту
		response := "Сообщение получено\n"
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Ошибка при отправке ответа:", err)
			return
		}

		// Проверяем, не завершился ли контекст (сигнал на завершение сервера)
		select {
		case <-ctx.Done():
			fmt.Println("Соединение завершено по сигналу остановки")
			return
		default:
		}
	}
}
