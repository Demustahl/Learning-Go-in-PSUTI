package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"log"
)

// Функция для создания хэшированного ключа фиксированной длины
func createHashKey(key string) []byte {
	// Приводим ключ к длине 32 байта для AES-256
	hashKey := make([]byte, 32)
	copy(hashKey, key)
	return hashKey
}

// Функция для шифрования данных
func encrypt(data, passphrase string) string {
	block, err := aes.NewCipher(createHashKey(passphrase))
	if err != nil {
		log.Fatalf("Ошибка при создании шифра: %v", err)
	}

	plainText := []byte(data)
	cfb := cipher.NewCFBEncrypter(block, createHashKey(passphrase)[:block.BlockSize()])
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return hex.EncodeToString(cipherText)
}

// Функция для расшифровки данных
func decrypt(data, passphrase string) string {
	block, err := aes.NewCipher(createHashKey(passphrase))
	if err != nil {
		log.Fatalf("Ошибка при создании шифра: %v", err)
	}

	cipherText, err := hex.DecodeString(data)
	if err != nil {
		log.Fatalf("Ошибка при декодировании шифртекста: %v", err)
	}

	cfb := cipher.NewCFBDecrypter(block, createHashKey(passphrase)[:block.BlockSize()])
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText)
}

func main() {
	fmt.Println("Выберите действие:")
	fmt.Println("1 - Зашифровать строку")
	fmt.Println("2 - Расшифровать строку")

	var choice int
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		fmt.Println("Введите строку для шифрования:")
		var input string
		fmt.Scanln(&input)

		fmt.Println("Введите секретный ключ:")
		var key string
		fmt.Scanln(&key)

		encrypted := encrypt(input, key)
		fmt.Printf("Зашифрованная строка: %s\n", encrypted)

	case 2:
		fmt.Println("Введите строку для расшифровки:")
		var input string
		fmt.Scanln(&input)

		fmt.Println("Введите секретный ключ:")
		var key string
		fmt.Scanln(&key)

		decrypted := decrypt(input, key)
		fmt.Printf("Расшифрованная строка: %s\n", decrypted)

	default:
		fmt.Println("Неверный выбор действия.")
	}
}
