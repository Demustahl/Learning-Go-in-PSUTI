package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Структура для обработки JSON-данных
type Data struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func main() {
	// Обработчик для GET-запроса на /hello
	http.HandleFunc("/hello", helloHandler)

	// Обработчик для POST-запроса на /data
	http.HandleFunc("/data", dataHandler)

	// Запуск сервера на порту 8080
	fmt.Println("HTTP-сервер запущен на порту 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}
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
