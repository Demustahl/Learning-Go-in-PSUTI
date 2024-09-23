package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Выберите задание:")
	fmt.Println("1 - Определить четность числа")
	fmt.Println("2 - Определить знак числа")
	fmt.Println("3 - Вывести числа от 1 до 10")
	fmt.Println("4 - Найти длину строки")
	fmt.Println("5 - Вычислить площадь прямоугольника")
	fmt.Println("6 - Найти среднее значение двух чисел")

	var choice int
	fmt.Print("Введите номер: ")
	fmt.Scanln(&choice)
	fmt.Println()

	switch choice {
	case 1:
		checkEvenOdd()
	case 2:
		checkSign()
	case 3:
		printNumbers()
	case 4:
		getStringLength()
	case 5:
		calculateRectangleArea()
	case 6:
		calculateAverage()
	default:
		fmt.Println("Неверный выбор")
	}
}

// Проверка четности числа
func checkEvenOdd() {
	var num int
	fmt.Println("Введите число:")
	fmt.Scanln(&num)

	if num%2 == 0 {
		fmt.Println("Четное")
	} else {
		fmt.Println("Нечетное")
	}
}

// Определение знака числа
func checkSign() {
	var num int
	fmt.Println("Введите число:")
	fmt.Scanln(&num)

	switch {
	case num > 0:
		fmt.Println("Положительное")
	case num < 0:
		fmt.Println("Отрицательное")
	default:
		fmt.Println("Это ноль")
	}
}

// Вывода чисел от 1 до 10
func printNumbers() {
	for i := 1; i <= 10; i++ {
		fmt.Println(i)
	}
}

// Подсчет длины строки
func getStringLength() {
	// var input string
	// fmt.Println("Введите строку:")
	// fmt.Scanln(&input)  // считывает до пробела
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Введите строку:")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input) // убирает лишние пробелы и символ новой строки

	fmt.Println("Длина строки:", len(input))
}

// Определение структуры
type Rectangle struct {
	Width, Height float64
}

// Метод для вычисления площади
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Создание прямоугольника и подсчёт его площади
func calculateRectangleArea() {
	var width, height float64
	fmt.Println("Введите ширину и высоту прямоугольника:")
	fmt.Scanln(&width, &height)

	rect := Rectangle{Width: width, Height: height}
	fmt.Println("Площадь прямоугольника:", rect.Area())
}

// Нахождение среднего значения двух чисел
func calculateAverage() {
	var num1, num2 int
	fmt.Println("Введите два числа:")
	fmt.Scanln(&num1, &num2)

	average := float64(num1+num2) / 2
	fmt.Println("Среднее значение:", average)
}
