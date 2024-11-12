// client.go
package main

import (
	"bufio"
	"crypto"
	"crypto/rand"
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

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	privateKeyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	data, _ := pem.Decode(privateKeyData)
	if data == nil || data.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("Неверный формат приватного ключа")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(data.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func main() {
	// Загружаем приватный ключ клиента
	privateKey, err := loadPrivateKey("client_private_key.pem")
	if err != nil {
		log.Fatalf("Ошибка при загрузке приватного ключа: %v", err)
	}

	// Подключаемся к серверу
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Ошибка при подключении к серверу: %v", err)
	}
	defer conn.Close()

	// Читаем сообщение для отправки
	fmt.Println("Введите сообщение для отправки:")
	reader := bufio.NewReader(os.Stdin)
	message, _ := reader.ReadString('\n')
	message = message[:len(message)-1] // удаляем символ новой строки

	// Подписываем сообщение
	hashed := sha256.Sum256([]byte(message))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		log.Fatalf("Ошибка при подписании сообщения: %v", err)
	}

	// Отправляем сообщение и подпись серверу
	encoder := gob.NewEncoder(conn)
	err = encoder.Encode(message)
	if err != nil {
		log.Fatalf("Ошибка при отправке сообщения: %v", err)
	}
	err = encoder.Encode(signature)
	if err != nil {
		log.Fatalf("Ошибка при отправке подписи: %v", err)
	}

	// Получаем ответ от сервера
	serverReader := bufio.NewReader(conn)
	response, err := serverReader.ReadString('\n')
	if err != nil {
		log.Fatalf("Ошибка при получении ответа от сервера: %v", err)
	}

	fmt.Println("Ответ от сервера:", response)
}
