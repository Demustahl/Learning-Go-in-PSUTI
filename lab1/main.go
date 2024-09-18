package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Выберите задание:")
	fmt.Println("1 - Вывести текущее время и дату")
	fmt.Println("2 - Создать и вывести переменные различных типов")
	fmt.Println("3 - Выполнить арифметические операции с целыми числами")
	fmt.Println("4 - Вычислить сумму и разность чисел с плавающей запятой")
	fmt.Println("5 - Вычислить среднее значение трех чисел")

	var choice int
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		printCurrentTime()
	case 2:
		printVariables()
	case 3:
		performArithmetic()
	case 4:
		calculateSumAndDifference()
	case 5:
		calculateAverage()
	default:
		fmt.Println("Неверный выбор")
	}

}

// Задание 1: вывод текущего времени и даты
func printCurrentTime() {
	currentTime := time.Now()
	fmt.Println("Текущая дата и время:", currentTime.Format("02-01-2006 15:04:05"))
}

// Задание 2 и 3: создание и вывод переменных различных типов
func printVariables() {
	var integer int = 42
	var floating float64 = 3.14
	var str string = "Пример строки"
	var boolean bool = true

	// Задание 3: использование краткой формы объявления переменных
	x := 100
	y := 12.34
	z := "Краткая форма"
	w := false

	fmt.Println("Целое число:", integer)
	fmt.Println("Число с плавающей запятой:", floating)
	fmt.Println("Строка:", str)
	fmt.Println("Булево значение:", boolean)

	fmt.Println("Переменные, объявленные краткой формой:")
	fmt.Println("x:", x)
	fmt.Println("y:", y)
	fmt.Println("z:", z)
	fmt.Println("w:", w)
}

// Задание 3: арифметические операции с целыми числами
func performArithmetic() {
	var a, b int
	fmt.Println("Введите два целых числа:")
	fmt.Scanln(&a, &b)

	fmt.Println("Сумма:", a+b)
	fmt.Println("Разность:", a-b)
	fmt.Println("Произведение:", a*b)
	if b != 0 {
		fmt.Println("Частное:", a/b)
	} else {
		fmt.Println("Деление на ноль!")
	}
}

// Задание 4: сумма и разность чисел с плавающей запятой
func calculateSumAndDifference() {
	var x, y float64
	fmt.Println("Введите два числа с плавающей запятой (через пробел):")
	fmt.Scanln(&x, &y)

	fmt.Printf("Сумма: %f\n", x+y)
	fmt.Printf("Разность: %f\n", x-y)
}

// Задание 5: вычисление среднего значения трех чисел
func calculateAverage() {
	var a, b, c float64
	fmt.Println("Введите три числа для вычисления среднего значения (через пробел):")
	fmt.Scanln(&a, &b, &c)

	average := (a + b + c) / 3
	fmt.Printf("Среднее значение: %f\n", average)
}
