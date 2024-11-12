// server.go
package main

import (
	"bufio"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/gob"
	"encoding/pem"
	"fmt"
	"log"
	"net"
	"os"
)

func loadPublicKey(path string) (*rsa.PublicKey, error) {
	publicKeyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	data, _ := pem.Decode(publicKeyData)
	if data == nil || data.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("Неверный формат публичного ключа")
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(data.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("Неверный тип публичного ключа")
	}

	return publicKey, nil
}

func handleConnection(conn net.Conn, clientPublicKey *rsa.PublicKey) {
	defer conn.Close()
	decoder := gob.NewDecoder(conn)

	var message string
	var signature []byte

	// Получаем сообщение
	err := decoder.Decode(&message)
	if err != nil {
		log.Println("Ошибка при получении сообщения:", err)
		return
	}

	// Получаем подпись
	err = decoder.Decode(&signature)
	if err != nil {
		log.Println("Ошибка при получении подписи:", err)
		return
	}

	fmt.Println("Получено сообщение:", message)

	// Проверяем подпись
	hashed := sha256.Sum256([]byte(message))
	err = rsa.VerifyPKCS1v15(clientPublicKey, crypto.SHA256, hashed[:], signature)
	var response string
	if err != nil {
		fmt.Println("Подпись недействительна.")
		response = "Подпись недействительна."
	} else {
		fmt.Println("Подпись подтверждена.")
		response = "Подпись подтверждена."
	}

	// Отправляем ответ клиенту
	writer := bufio.NewWriter(conn)
	writer.WriteString(response + "\n")
	writer.Flush()
}

func main() {
	// Загружаем публичный ключ клиента
	clientPublicKey, err := loadPublicKey("client_public_key.pem")
	if err != nil {
		log.Fatalf("Ошибка при загрузке публичного ключа клиента: %v", err)
	}

	// Слушаем на порту 8080
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
	defer ln.Close()

	fmt.Println("Сервер запущен и ожидает подключения...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Ошибка при подключении клиента:", err)
			continue
		}

		go handleConnection(conn, clientPublicKey)
	}
}
