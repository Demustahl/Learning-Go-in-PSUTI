package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Структура для обработки JSON-данных
type Data struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func main() {
	// Создаем маршрутизатор
	mux := http.NewServeMux()

	// Добавляем обработчики с middleware
	mux.Handle("/hello", loggingMiddleware(http.HandlerFunc(helloHandler)))
	mux.Handle("/data", loggingMiddleware(http.HandlerFunc(dataHandler)))

	// Запуск сервера на порту 8080 с использованием маршрутизатора
	fmt.Println("HTTP-сервер запущен на порту 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}

// Middleware для логирования
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Логируем метод и URL запроса
		log.Printf("Начало обработки запроса: %s %s", r.Method, r.URL.Path)

		// Передаем управление следующему обработчику
		next.ServeHTTP(w, r)

		// Логируем время выполнения запроса
		duration := time.Since(start)
		log.Printf("Запрос обработан: %s %s, время выполнения: %v", r.Method, r.URL.Path, duration)
	})
}

// Обработчик для GET /hello
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Отправляем приветственное сообщение
		fmt.Fprintf(w, "Hello, World!")
	} else {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// Обработчик для POST /data
func dataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Читаем тело запроса
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Ошибка при чтении данных", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// Декодируем JSON-данные
		var data Data
		err = json.Unmarshal(body, &data)
		if err != nil {
			http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
			return
		}

		// Выводим содержимое JSON в консоль
		fmt.Printf("Полученные данные: Name=%s, Value=%s\n", data.Name, data.Value)

		// Отправляем подтверждение клиенту
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Данные получены")
	} else {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}
