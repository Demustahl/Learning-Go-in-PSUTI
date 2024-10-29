package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Тестовые данные
var testUser = User{Name: "Test User", Age: 30}

// Создание тестового пользователя в MongoDB
func createTestUser() User {
	user := testUser
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		panic(err)
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return user
}

// Очистка базы данных после тестов
func clearDatabase() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userCollection.DeleteMany(ctx, bson.M{})
}

func TestMain(m *testing.M) {
	connectDB()
	clearDatabase()
	m.Run()
	clearDatabase()
}

// Тест функции getUsers
func TestGetUsers(t *testing.T) {
	req, _ := http.NewRequest("GET", "/users", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/users", getUsers).Methods("GET")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var users []User
	json.NewDecoder(rr.Body).Decode(&users)
	if len(users) != 0 {
		t.Errorf("Expected no users, got %v", len(users))
	}
}

// Тест функции createUser
func TestCreateUser(t *testing.T) {
	payload := `{"name": "Gandalf", "age": 100}`
	req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/users", createUser).Methods("POST")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var user User
	json.NewDecoder(rr.Body).Decode(&user)
	if user.Name != "Gandalf" || user.Age != 100 {
		t.Errorf("handler returned unexpected user data: got %v", user)
	}
}

// Тест функции getUser
func TestGetUser(t *testing.T) {
	createdUser := createTestUser()

	req, _ := http.NewRequest("GET", "/users/"+createdUser.ID.Hex(), nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", getUser).Methods("GET")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var user User
	json.NewDecoder(rr.Body).Decode(&user)
	if user.ID != createdUser.ID {
		t.Errorf("Expected user ID %v, got %v", createdUser.ID, user.ID)
	}
}

// Тест функции updateUser
func TestUpdateUser(t *testing.T) {
	createdUser := createTestUser()

	payload := `{"name": "Gandalf the White", "age": 101}`
	req, _ := http.NewRequest("PUT", "/users/"+createdUser.ID.Hex(), bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", updateUser).Methods("PUT")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var updatedUser User
	json.NewDecoder(rr.Body).Decode(&updatedUser)
	if updatedUser.Name != "Gandalf the White" || updatedUser.Age != 101 {
		t.Errorf("handler returned unexpected user data: got %v", updatedUser)
	}
}

// Тест функции deleteUser
func TestDeleteUser(t *testing.T) {
	createdUser := createTestUser()

	req, _ := http.NewRequest("DELETE", "/users/"+createdUser.ID.Hex(), nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	// Проверяем, что пользователя нет
	req, _ = http.NewRequest("GET", "/users/"+createdUser.ID.Hex(), nil)
	rr = httptest.NewRecorder()

	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Expected status %v, got %v", http.StatusNotFound, status)
	}
}

// Тест функции getUsers с пагинацией и фильтрацией
func TestGetUsersWithPaginationAndFilter(t *testing.T) {
	clearDatabase()
	createTestUser()
	createTestUser()

	req, _ := http.NewRequest("GET", "/users?page=1&pageSize=1", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/users", getUsers).Methods("GET")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var users []User
	json.NewDecoder(rr.Body).Decode(&users)
	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %v", len(users))
	}
}

// Тест фильтрации пользователей по имени
func TestGetUsersWithNameFilter(t *testing.T) {
	clearDatabase()
	createTestUser()

	req, _ := http.NewRequest("GET", "/users?name=Test", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/users", getUsers).Methods("GET")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var users []User
	json.NewDecoder(rr.Body).Decode(&users)
	if len(users) != 1 || users[0].Name != "Test User" {
		t.Errorf("Expected user with name 'Test User', got %v", users)
	}
}
