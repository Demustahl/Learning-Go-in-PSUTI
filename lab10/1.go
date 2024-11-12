package main

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"strings"
)

// Функция для выбора хэш-функции на основе пользовательского выбора
func getHashFunction(choice int) (hash.Hash, string) {
	switch choice {
	case 1:
		return md5.New(), "MD5"
	case 2:
		return sha256.New(), "SHA-256"
	case 3:
		return sha512.New(), "SHA-512"
	default:
		return nil, ""
	}
}

// Функция для хэширования строки с помощью выбранной хэш-функции
func hashString(input string, hasher hash.Hash) string {
	hasher.Write([]byte(input))        // Переводим строку в байты и хэшируем
	hashSum := hasher.Sum(nil)         // Получаем байтовое представление хэша
	return hex.EncodeToString(hashSum) // Переводим байты в строку в шестнадцатеричном формате
}

// Функция для проверки целостности строки с её хэшом
func verifyHash(input, providedHash string, hasher hash.Hash) bool {
	calculatedHash := hashString(input, hasher)            // Хэшируем введённую строку
	return strings.EqualFold(calculatedHash, providedHash) // Сравниваем, не учитывая регистр
}

func main() {
	fmt.Println("Выберите действие:")
	fmt.Println("1 - Захешировать строку")
	fmt.Println("2 - Проверить целостность строки с её хэшом")

	var action int
	fmt.Scan(&action)

	if action == 1 {
		// Запрос на хэширование строки
		fmt.Println("Введите строку для хэширования:")
		var input string
		fmt.Scan(&input)

		fmt.Println("Выберите алгоритм хэширования:")
		fmt.Println("1 - MD5")
		fmt.Println("2 - SHA-256")
		fmt.Println("3 - SHA-512")

		var choice int
		fmt.Scan(&choice)

		hasher, hashName := getHashFunction(choice)
		if hasher == nil {
			fmt.Println("Неверный выбор алгоритма.")
			return
		}

		hashResult := hashString(input, hasher)
		fmt.Printf("Алгоритм: %s\nХэш: %s\n", hashName, hashResult)

	} else if action == 2 {
		// Проверка целостности строки с хэшем
		fmt.Println("Введите строку для проверки:")
		var input string
		fmt.Scan(&input)

		fmt.Println("Введите предполагаемый хэш:")
		var providedHash string
		fmt.Scan(&providedHash)

		fmt.Println("Выберите алгоритм хэширования:")
		fmt.Println("1 - MD5")
		fmt.Println("2 - SHA-256")
		fmt.Println("3 - SHA-512")

		var choice int
		fmt.Scan(&choice)

		hasher, hashName := getHashFunction(choice)
		if hasher == nil {
			fmt.Println("Неверный выбор алгоритма.")
			return
		}

		if verifyHash(input, providedHash, hasher) {
			fmt.Printf("Алгоритм: %s\nХэш совпадает. Данные целостны.\n", hashName)
		} else {
			fmt.Printf("Алгоритм: %s\nХэш не совпадает. Данные изменены или хэш неверен.\n", hashName)
		}

	} else {
		fmt.Println("Неверный выбор действия.")
	}
}
