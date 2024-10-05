package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Глобальная переменная people
var people = map[string]int{
	"John":      25,
	"Bob":       22,
	"Frodo":     33,
	"SpongeBob": 20,
	"Mickey":    28,
	"Yoda":      900,
	"Dobby":     18,
	"Hobbiton":  21,
}

func main() {
	fmt.Println("Выберите задание для выполнения:")
	fmt.Println("1 - Добавить и вывести записи о людях")
	fmt.Println("2 - Вывести средний возраст людей")
	fmt.Println("3 - Удалить запись по имени")
	fmt.Println("4 - Ввести строку и вывести её в верхнем регистре")
	fmt.Println("5 - Ввести числа и вывести их сумму")
	fmt.Println("6 - Ввести массив чисел и вывести его в обратном порядке")

	var task int
	fmt.Scanln(&task)

	switch task {
	case 1:
		addAndDisplayPeople()
	case 2:
		calculateAndDisplayAverageAge()
	case 3:
		deletePerson()
	case 4:
		toUpperCase()
	case 5:
		sumNumbers()
	case 6:
		reverseArray()
	default:
		fmt.Println("Неверный выбор. Попробуйте снова.")
	}
}

// 1. Добавить и вывести записи о людях
func addAndDisplayPeople() {
	// Выводим записи до добавления
	fmt.Println("Записи до добавления нового человека:")
	for name, age := range people {
		fmt.Printf("%s: %d лет\n", name, age)
	}

	// Добавляем запись о Чубакке
	people["Chewbacca"] = 200

	// Чтение новой записи от пользователя
	fmt.Println("\nВведите имя нового человека:")
	reader := bufio.NewReader(os.Stdin)
	newName, _ := reader.ReadString('\n')
	newName = strings.TrimSpace(newName)

	fmt.Println("Введите возраст нового человека:")
	var newAge int
	fmt.Scanln(&newAge)

	// Добавляем новую запись
	people[newName] = newAge

	// Выводим записи после добавления
	fmt.Println("\nЗаписи после добавления нового человека:")
	for name, age := range people {
		fmt.Printf("%s: %d лет\n", name, age)
	}
}

// 2. Вычисление и вывод среднего возраста
func calculateAndDisplayAverageAge() {
	averageAge := calculateAverageAge(people)
	fmt.Printf("\nСредний возраст: %.2f\n", averageAge)
}

// Функция для вычисления среднего возраста
func calculateAverageAge(people map[string]int) float64 {
	var sum, count int
	for _, age := range people {
		sum += age
		count++
	}
	return float64(sum) / float64(count)
}

// 3. Удаление записи по имени
func deletePerson() {
	// Выводим записи до удаления
	fmt.Println("Записи до удаления:")
	for name, age := range people {
		fmt.Printf("%s: %d лет\n", name, age)
	}

	// Удаление записи
	fmt.Println("\nВведите имя для удаления:")
	var nameToDelete string
	fmt.Scanln(&nameToDelete)
	delete(people, nameToDelete)

	// Выводим записи после удаления
	fmt.Println("\nЗаписи после удаления:")
	for name, age := range people {
		fmt.Printf("%s: %d лет\n", name, age)
	}
}

// 4. Чтение строки и вывод в верхнем регистре
func toUpperCase() {
	fmt.Println("Введите строку:")
	reader := bufio.NewReader(os.Stdin)
	inputString, _ := reader.ReadString('\n')
	fmt.Printf("Строка в верхнем регистре: %s\n", strings.ToUpper(inputString))
}

// 5. Чтение чисел и вывод их суммы
func sumNumbers() {
	fmt.Println("Введите числа через пробел:")
	reader := bufio.NewReader(os.Stdin)
	inputNumbers, _ := reader.ReadString('\n')
	numberStrings := strings.Fields(inputNumbers)
	var sum int
	for _, numStr := range numberStrings {
		num, _ := strconv.Atoi(numStr)
		sum += num
	}
	fmt.Printf("Сумма чисел: %d\n", sum)
}

// 6. Чтение массива чисел и вывод в обратном порядке
func reverseArray() {
	fmt.Println("Введите массив целых чисел через пробел:")
	reader := bufio.NewReader(os.Stdin)
	inputArray, _ := reader.ReadString('\n')
	numberArray := strings.Fields(inputArray)
	fmt.Println("Массив в обратном порядке:")
	for i := len(numberArray) - 1; i >= 0; i-- {
		fmt.Printf("%s ", numberArray[i])
	}
	fmt.Println()
}
