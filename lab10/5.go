// main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("my_secret_key")

// Структура пользователя
type User struct {
	Username string
	Password string
	Role     string
}

// Фиксированные данные пользователей для примера
var users = map[string]User{
	"admin": {"admin", "admin123", "admin"},
	"user":  {"user", "user123", "user"},
}

// Структура для JWT Claims
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// Функция для генерации JWT-токена
func GenerateJWT(username, role string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Функция для проверки JWT-токена и извлечения Claims
func ValidateJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("Invalid token")
	}

	return claims, nil
}

// Обработчик для входа пользователя
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, exists := users[username]
	if !exists || user.Password != password {
		http.Error(w, "Неверные имя пользователя или пароль", http.StatusUnauthorized)
		return
	}

	tokenString, err := GenerateJWT(user.Username, user.Role)
	if err != nil {
		http.Error(w, "Ошибка при генерации токена", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
	})

	fmt.Fprintf(w, "Успешный вход. Токен сохранён в cookie.")
}

// Обработчик для выхода пользователя
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	fmt.Fprintf(w, "Успешный выход из системы.")
}

// Обработчик для защищённого маршрута
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*Claims)
	if !ok {
		http.Error(w, "Не удалось получить данные пользователя", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Привет, %s! Это защищённый маршрут.", claims.Username)
}

// Обработчик для маршрута только для админа
func AdminHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*Claims)
	if !ok {
		http.Error(w, "Не удалось получить данные пользователя", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Привет, админ %s! У тебя есть доступ к этой информации.", claims.Username)
}

// Middleware для аутентификации
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Необходима аутентификация", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Необходима аутентификация", http.StatusBadRequest)
			return
		}

		tokenStr := cookie.Value
		claims, err := ValidateJWT(tokenStr)
		if err != nil {
			http.Error(w, "Неверный или просроченный токен", http.StatusUnauthorized)
			return
		}

		// Добавляем claims в контекст запроса
		ctx := r.Context()
		ctx = context.WithValue(ctx, "claims", claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// Middleware для проверки роли пользователя
func RoleMiddleware(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value("claims").(*Claims)
			if !ok {
				http.Error(w, "Не удалось получить данные пользователя", http.StatusInternalServerError)
				return
			}

			if claims.Role != role {
				http.Error(w, "Доступ запрещён: недостаточно прав", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/logout", LogoutHandler)

	// Защищённые маршруты
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/protected", ProtectedHandler)
	protectedMux.HandleFunc("/admin", AdminHandler)

	// Применение middleware для аутентификации
	http.Handle("/protected", AuthMiddleware(protectedMux))
	http.Handle("/admin", AuthMiddleware(RoleMiddleware("admin")(protectedMux)))

	fmt.Println("Сервер запущен на порту 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
