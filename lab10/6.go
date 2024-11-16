// main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Инициализация хранилища сессий
var store = sessions.NewCookieStore([]byte("super-secret-key"))

// Данные пользователей
var users = map[string]string{
	"admin": "admin123",
	"user":  "user123",
}

// Обработчик для главной страницы
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Добро пожаловать на главную страницу!")
}

// Обработчик для входа пользователя
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Устанавливаем Content-Type с указанием кодировки UTF-8
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// Отображение формы входа с CSRF-токеном
		token := csrf.Token(r)
		fmt.Fprintf(w, `
            <form action="/login" method="POST">
                <input type="hidden" name="gorilla.csrf.Token" value="%s">
                Имя пользователя: <input type="text" name="username"><br>
                Пароль: <input type="password" name="password"><br>
                <input type="submit" value="Войти">
            </form>
        `, token)
		return
	}

	// Обработка POST-запроса для входа
	username := r.FormValue("username")
	password := r.FormValue("password")

	if pwd, ok := users[username]; ok && pwd == password {
		session, _ := store.Get(r, "session-id")
		session.Values["authenticated"] = true
		session.Values["username"] = username
		session.Save(r, w)
		http.Redirect(w, r, "/protected", http.StatusFound)
	} else {
		http.Error(w, "Неверные имя пользователя или пароль", http.StatusUnauthorized)
	}
}

// Обработчик для выхода пользователя
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-id")
	session.Values["authenticated"] = false
	delete(session.Values, "username")
	session.Save(r, w)
	fmt.Fprintf(w, "Вы вышли из системы.")
}

// Middleware для проверки аутентификации
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-id")
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "Доступ запрещён. Необходимо войти в систему.", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Обработчик для защищённого маршрута
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-id")
	username, _ := session.Values["username"].(string)
	fmt.Fprintf(w, "Привет, %s! Это защищённый маршрут.", username)
}

func main() {
	r := mux.NewRouter()

	// Инициализация защиты от CSRF
	csrfMiddleware := csrf.Protect(
		[]byte("32-byte-long-auth-key-1234567890abcdef"),
		csrf.Secure(false), // Для разработки, в продакшене должно быть true
	)

	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/login", LoginHandler)
	r.HandleFunc("/logout", LogoutHandler)
	r.Handle("/protected", AuthMiddleware(http.HandlerFunc(ProtectedHandler)))

	// Применяем CSRF-защиту ко всем маршрутам
	http.Handle("/", csrfMiddleware(r))

	fmt.Println("Сервер запущен на порту 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
