package main

import (
	"fmt"
	"lab3/mathutils"
	"lab3/stringutils"
)

func main() {
	fmt.Println("Выберите задание для выполнения:")
	fmt.Println("1 - Найти факториал числа")
	fmt.Println("2 - Перевернуть строку")
	fmt.Println("3 - Показать список из 5 чисел")
	fmt.Println("4 - Добавить и удалить элементы в срезе")
	fmt.Println("5 - Найти самую длинную строку")

	var task int
	fmt.Scanln(&task)

	switch task {
	case 1:
		getFactorial()
	case 2:
		flipString()
	case 3:
		displayIntArray()
	case 4:
		modifySlice()
	case 5:
		longestWord()
	default:
		fmt.Println("Неверный выбор. Попробуйте снова.")
	}
}

// Задание 1: Найти факториал числа
func getFactorial() {
	var number int
	fmt.Print("Введите целое число: ")
	fmt.Scanln(&number)
	output := mathutils.Factorial(number)
	fmt.Printf("Факториал числа %d равен %d\n", number, output)
}

// Задание 2: Перевернуть строку
func flipString() {
	var input string
	fmt.Print("Введите строку: ")
	fmt.Scanln(&input)
	output := stringutils.Reverse(input)
	fmt.Printf("Перевернутая строка: %s\n", output)
}

// Задание 3: Показать список из 5 чисел
func displayIntArray() {
	nums := [5]int{13, 4, 69, 404, 666}
	fmt.Println("Список чисел:", nums)
}

// Задание 4: Добавить и удалить элементы в срезе
func modifySlice() {
	origArray := [5]int{1, -5, 13455, 88, -111}
	partSlice := origArray[2:4] // Выбираем часть массива
	fmt.Println("Исходный срез:", partSlice)

	partSlice = append(partSlice, 33) // Добавляем элемент
	fmt.Println("Срез после добавления:", partSlice)

	partSlice = partSlice[:len(partSlice)-1] // Удаляем последний элемент
	fmt.Println("Срез после удаления:", partSlice)
}

// Задание 5: Найти самую длинную строку
func longestWord() {
	words := []string{".", "Moscow", "W_W", "1234567890"}
	var longestWord string

	for _, word := range words {
		if len(word) > len(longestWord) {
			longestWord = word
		}
	}

	fmt.Printf("Самая длинная строка: %s\n", longestWord)
}
