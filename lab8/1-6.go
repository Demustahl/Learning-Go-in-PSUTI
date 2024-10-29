package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
	Age  int                `json:"age" bson:"age"`
}

var client *mongo.Client
var userCollection *mongo.Collection

func connectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Устанавливаем соединение с MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Проверяем соединение
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Успешно подключились к MongoDB!")

	// Получаем ссылку на коллекцию
	userCollection = client.Database("dbLab8").Collection("users")
}

// getUsers с поддержкой пагинации и фильтрации
func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Чтение параметров запроса для пагинации
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize < 1 {
		pageSize = 10 // значение по умолчанию
	}

	// Чтение параметров запроса для фильтрации
	filter := bson.M{}
	name := r.URL.Query().Get("name")
	if name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"} // поиск по имени (регистр игнорируется)
	}

	ageStr := r.URL.Query().Get("age")
	if ageStr != "" {
		age, err := strconv.Atoi(ageStr)
		if err == nil {
			filter["age"] = age // фильтрация по точному значению возраста
		}
	}

	// Параметры для запроса к MongoDB
	findOptions := options.Find()
	findOptions.SetSkip(int64((page - 1) * pageSize))
	findOptions.SetLimit(int64(pageSize))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Выполняем запрос к MongoDB
	cursor, err := userCollection.Find(ctx, filter, findOptions)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve users")
		return
	}
	defer cursor.Close(ctx)

	// Получение данных пользователей
	var users []User
	for cursor.Next(ctx) {
		var user User
		cursor.Decode(&user)
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Error iterating through users")
		return
	}

	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = userCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	// Валидация данных пользователя
	isValid, validationErr := validateUser(user)
	if !isValid {
		sendErrorResponse(w, http.StatusBadRequest, validationErr)
		return
	}

	user.ID = primitive.NewObjectID()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var user User
	json.NewDecoder(r.Body).Decode(&user)

	// Валидация данных пользователя
	isValid, validationErr := validateUser(user)
	if !isValid {
		sendErrorResponse(w, http.StatusBadRequest, validationErr)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := userCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"name": user.Name, "age": user.Age}},
	)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to update user")
		return
	}
	if result.MatchedCount == 0 {
		sendErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	user.ID = id
	json.NewEncoder(w).Encode(user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := userCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if result.DeletedCount == 0 {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func sendErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// validateUser проверяет, соответствуют ли данные пользователя нашим условиям.
func validateUser(user User) (bool, string) {
	if strings.TrimSpace(user.Name) == "" {
		return false, "Name cannot be empty"
	}
	if user.Age <= 0 {
		return false, "Age must be greater than zero"
	}
	return true, ""
}

func main() {
	connectDB()

	router := mux.NewRouter()
	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
