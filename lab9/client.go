package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const serverURL = "http://localhost:8080"

type User struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var sessionToken string

func clearConsole() {
	fmt.Print("\033[H\033[2J")
	// fmt.Print("") // это просто, чтобы убрать временно очистку экрана
}

func main() {
	for {
		clearConsole()
		fmt.Println("Выберите операцию:")

		if sessionToken == "" {
			fmt.Println("1. Регистрация")
			fmt.Println("2. Вход")
		}

		fmt.Println("3. Добавить пользователя")
		fmt.Println("4. Прочитать информацию о пользователе")
		fmt.Println("5. Обновить информацию о пользователе")
		fmt.Println("6. Удалить пользователя")
		fmt.Println("7. Вывести список пользователей")
		if sessionToken != "" {
			fmt.Println("9. Выход из аккаунта")
		}
		fmt.Println("8. Выход")

		var choice int
		fmt.Scanln(&choice)

		clearConsole()

		if sessionToken == "" {
			switch choice {
			case 1:
				registerUser()
			case 2:
				loginUser()
			case 8:
				fmt.Println("Выход из программы")
				os.Exit(0)
			default:
				fmt.Println("Вы не авторизованы. Пожалуйста, войдите в систему.")
			}
		} else {
			switch choice {
			case 3:
				createUser()
			case 4:
				readUser()
			case 5:
				updateUser()
			case 6:
				deleteUser()
			case 7:
				listUsers()
			case 9:
				logoutUser()
			case 8:
				// Добавим вызов logoutUser() перед выходом
				logoutUser()
				fmt.Println("Выход из программы")
				os.Exit(0)
			default:
				fmt.Println("Неверный выбор. Попробуйте еще раз.")
			}
		}

		fmt.Println("\nНажмите Enter для продолжения...")
		fmt.Scanln()
	}
}

func logoutUser() {
	req, err := http.NewRequest(http.MethodPost, serverURL+"/logout", nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}
	req.Header.Set("Authorization", sessionToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Ошибка:", result["error"])
	} else {
		fmt.Println(result["message"])
		sessionToken = ""
	}
}

func registerUser() {
	fmt.Println("Регистрация нового пользователя")
	fmt.Println("-------------------------------")

	var user User
	fmt.Print("Введите имя пользователя: ")
	fmt.Scanln(&user.Name)
	fmt.Print("Введите email пользователя: ")
	fmt.Scanln(&user.Email)
	fmt.Print("Введите пароль: ")
	fmt.Scanln(&user.Password)

	jsonData, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Ошибка при кодировании данных:", err)
		return
	}

	resp, err := http.Post(serverURL+"/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if resp.StatusCode != http.StatusCreated {
		fmt.Println("Ошибка:", result["error"])
	} else {
		fmt.Println(result["message"])
		fmt.Printf("Ваш ID пользователя: %v\n", result["userID"])
	}
}

func loginUser() {
	fmt.Println("Вход пользователя")
	fmt.Println("-----------------")

	var creds struct {
		ID       int    `json:"id"`
		Password string `json:"password"`
	}

	fmt.Print("Введите ID пользователя: ")
	fmt.Scanln(&creds.ID)
	fmt.Print("Введите пароль: ")
	fmt.Scanln(&creds.Password)

	jsonData, err := json.Marshal(creds)
	if err != nil {
		fmt.Println("Ошибка при кодировании данных:", err)
		return
	}

	resp, err := http.Post(serverURL+"/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Ошибка:", result["error"])
	} else {
		fmt.Println(result["message"])
		sessionToken = result["token"].(string)
		fmt.Println("Токен сессии сохранен.")
	}
}

func createUser() {
	if sessionToken == "" {
		fmt.Println("Вы не авторизованы. Пожалуйста, войдите в систему.")
		return
	}

	fmt.Println("Добавление нового пользователя")
	fmt.Println("------------------------------")

	var user User
	fmt.Print("Введите имя пользователя: ")
	fmt.Scanln(&user.Name)
	fmt.Print("Введите email пользователя: ")
	fmt.Scanln(&user.Email)
	fmt.Print("Введите пароль: ")
	fmt.Scanln(&user.Password)

	jsonData, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Ошибка при кодировании данных:", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, serverURL+"/user", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", sessionToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if resp.StatusCode != http.StatusCreated {
		fmt.Println("Ошибка:", result["error"])
	} else {
		fmt.Println(result["message"])
		fmt.Printf("ID созданного пользователя: %v\n", result["userID"])
	}
}

func readUser() {
	if sessionToken == "" {
		fmt.Println("Вы не авторизованы. Пожалуйста, войдите в систему.")
		return
	}

	fmt.Println("Просмотр информации о пользователе")
	fmt.Println("----------------------------------")

	var id int
	fmt.Print("Введите ID пользователя: ")
	fmt.Scanln(&id)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/user/%d", serverURL, id), nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}
	req.Header.Set("Authorization", sessionToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResult map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResult)
		fmt.Println("Ошибка:", errResult["error"])
		return
	}

	var user User
	json.NewDecoder(resp.Body).Decode(&user)
	fmt.Println("Информация о пользователе:")
	fmt.Println("--------------------------")
	fmt.Printf("ID: %d\nИмя: %s\nEmail: %s\n", user.ID, user.Name, user.Email)
}

func updateUser() {
	if sessionToken == "" {
		fmt.Println("Вы не авторизованы. Пожалуйста, войдите в систему.")
		return
	}

	fmt.Println("Обновление информации о пользователе")
	fmt.Println("------------------------------------")

	var id int
	fmt.Print("Введите ID пользователя для обновления: ")
	fmt.Scanln(&id)

	var user User
	fmt.Print("Введите новое имя пользователя: ")
	fmt.Scanln(&user.Name)
	fmt.Print("Введите новый email пользователя: ")
	fmt.Scanln(&user.Email)
	fmt.Print("Введите новый пароль: ")
	fmt.Scanln(&user.Password)
	user.ID = id

	jsonData, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Ошибка при кодировании данных:", err)
		return
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/user/%d", serverURL, id), bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", sessionToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Ошибка:", result["error"])
	} else {
		fmt.Println(result["message"])
	}
}

func deleteUser() {
	if sessionToken == "" {
		fmt.Println("Вы не авторизованы. Пожалуйста, войдите в систему.")
		return
	}

	fmt.Println("Удаление пользователя")
	fmt.Println("----------------------")

	var id int
	fmt.Print("Введите ID пользователя для удаления: ")
	fmt.Scanln(&id)

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/user/%d", serverURL, id), nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}
	req.Header.Set("Authorization", sessionToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Ошибка:", result["error"])
	} else {
		fmt.Println(result["message"])
	}
}

func listUsers() {
	if sessionToken == "" {
		fmt.Println("Вы не авторизованы. Пожалуйста, войдите в систему.")
		return
	}

	fmt.Println("Список пользователей")
	fmt.Println("--------------------")

	req, err := http.NewRequest(http.MethodGet, serverURL+"/users", nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}
	req.Header.Set("Authorization", sessionToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResult map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResult)
		fmt.Println("Ошибка:", errResult["error"])
		return
	}

	var users []User
	json.NewDecoder(resp.Body).Decode(&users)

	if len(users) == 0 {
		fmt.Println("Список пользователей пуст.")
		return
	}

	for _, user := range users {
		fmt.Printf("ID: %d\nИмя: %s\nEmail: %s\n", user.ID, user.Name, user.Email)
		fmt.Println("--------------------")
	}
}
