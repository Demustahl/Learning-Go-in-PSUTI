package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

var (
	users       = make(map[int]User)
	sessions    = make(map[string]int) // sessionToken -> userID
	activeUsers = make(map[int]bool)   // userID -> isActive
	mu          sync.Mutex
)

func main() {
	// Загрузка пользователей из файла при запуске сервера
	err := loadUsers()
	if err != nil {
		log.Fatalf("Ошибка загрузки данных пользователей: %v", err)
	}

	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", authMiddleware(logoutHandler))
	http.HandleFunc("/user", authMiddleware(createUserHandler))
	http.HandleFunc("/user/", authMiddleware(userHandler))
	http.HandleFunc("/users", authMiddleware(listUsersHandler))

	fmt.Println("Сервер запущен на порту 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Функция загрузки пользователей из файла
func loadUsers() error {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.Open("users.json")
	if err != nil {
		if os.IsNotExist(err) {
			users = make(map[int]User)
			return nil
		}
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&users)
	if err != nil {
		return err
	}
	return nil
}

// Функция сохранения пользователей в файл
func saveUsers() error {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.Create("users.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(users)
	if err != nil {
		return err
	}
	return nil
}

func generateUserID() int {
	usedIDs := make(map[int]bool)
	for id := range users {
		usedIDs[id] = true
	}
	// Ищем минимально доступный ID, начиная с 1
	for i := 1; ; i++ {
		if !usedIDs[i] {
			return i
		}
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректные данные")
		return
	}

	mu.Lock()
	user.ID = generateUserID()
	mu.Unlock()

	// Хешируем пароль перед сохранением
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка при обработке пароля")
		return
	}
	user.Password = string(hashedPassword)

	mu.Lock()
	users[user.ID] = user
	mu.Unlock()

	// Сохранение пользователей в файл
	err = saveUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка сохранения данных")
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Пользователь успешно зарегистрирован",
		"userID":  user.ID,
	})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
		return
	}

	var creds struct {
		ID       int    `json:"id"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректные данные")
		return
	}

	mu.Lock()
	user, exists := users[creds.ID]
	if !exists {
		mu.Unlock()
		respondWithError(w, http.StatusUnauthorized, "Неверный ID пользователя или пароль")
		return
	}

	// Проверяем, авторизован ли пользователь
	if activeUsers[creds.ID] {
		mu.Unlock()
		respondWithError(w, http.StatusUnauthorized, "Пользователь уже авторизован в системе")
		return
	}
	mu.Unlock()

	// Сравниваем хешированный пароль
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Неверный ID пользователя или пароль")
		return
	}

	token := generateSessionToken()

	mu.Lock()
	sessions[token] = user.ID
	activeUsers[user.ID] = true // Отмечаем пользователя как активного
	mu.Unlock()

	respondWithJSON(w, http.StatusOK, map[string]interface{}{"message": "Успешный вход", "token": token})
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
		return
	}

	token := r.Header.Get("Authorization")
	if token == "" {
		respondWithError(w, http.StatusUnauthorized, "Токен не предоставлен")
		return
	}

	mu.Lock()
	userID, exists := sessions[token]
	if exists {
		delete(sessions, token)
		delete(activeUsers, userID)
	}
	mu.Unlock()

	if !exists {
		respondWithError(w, http.StatusUnauthorized, "Недействительный токен")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Вы успешно вышли из аккаунта"})
}

func generateSessionToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			respondWithError(w, http.StatusUnauthorized, "Токен не предоставлен")
			return
		}

		mu.Lock()
		userID, exists := sessions[token]
		mu.Unlock()

		if !exists {
			respondWithError(w, http.StatusUnauthorized, "Недействительный токен")
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		next(w, r.WithContext(ctx))
	}
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректные данные")
		return
	}

	mu.Lock()
	user.ID = generateUserID()
	mu.Unlock()

	// Хешируем пароль перед сохранением
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка при обработке пароля")
		return
	}
	user.Password = string(hashedPassword)

	mu.Lock()
	users[user.ID] = user
	mu.Unlock()

	// Сохранение пользователей в файл
	err = saveUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка сохранения данных")
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Пользователь создан",
		"userID":  user.ID,
	})
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		readUserHandler(w, r)
	case http.MethodPut:
		updateUserHandler(w, r)
	case http.MethodDelete:
		deleteUserHandler(w, r)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
	}
}

func readUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/user/"):]
	var id int
	fmt.Sscanf(idStr, "%d", &id)

	mu.Lock()
	user, exists := users[id]
	mu.Unlock()

	if !exists {
		respondWithError(w, http.StatusNotFound, "Пользователь не найден")
		return
	}

	// Не отправляем пароль клиенту
	user.Password = ""
	respondWithJSON(w, http.StatusOK, user)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		respondWithError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
		return
	}

	idStr := r.URL.Path[len("/user/"):]
	var id int
	fmt.Sscanf(idStr, "%d", &id)

	var updatedData User
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректные данные")
		return
	}

	mu.Lock()
	user, exists := users[id]
	mu.Unlock()

	if !exists {
		respondWithError(w, http.StatusNotFound, "Пользователь не найден")
		return
	}

	// Обновляем данные пользователя
	user.Name = updatedData.Name
	user.Email = updatedData.Email

	// Если предоставлен новый пароль, хешируем его
	if updatedData.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedData.Password), bcrypt.DefaultCost)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Ошибка при обработке пароля")
			return
		}
		user.Password = string(hashedPassword)
	}

	mu.Lock()
	users[id] = user
	mu.Unlock()

	// Сохранение пользователей в файл
	err := saveUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка сохранения данных")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Информация о пользователе обновлена"})
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondWithError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
		return
	}

	idStr := r.URL.Path[len("/user/"):]
	var id int
	fmt.Sscanf(idStr, "%d", &id)

	mu.Lock()
	_, exists := users[id]
	if exists {
		delete(users, id)
	}
	mu.Unlock()

	if !exists {
		respondWithError(w, http.StatusNotFound, "Пользователь не найден")
		return
	}

	// Сохранение пользователей в файл
	err := saveUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка сохранения данных")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Пользователь удален"})
}

func listUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
		return
	}

	mu.Lock()
	var userList []User
	for _, user := range users {
		user.Password = "" // Не отправляем пароли
		userList = append(userList, user)
	}
	mu.Unlock()

	respondWithJSON(w, http.StatusOK, userList)
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка кодирования JSON")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJSON(w, status, map[string]string{"error": message})
}
