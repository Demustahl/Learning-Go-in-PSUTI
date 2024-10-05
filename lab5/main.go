package main

import (
	"fmt"
	"math"
)

// Структура Person
type Person struct {
	name string
	age  int
}

// Метод для вывода информации о человеке
func (p Person) Info() {
	fmt.Printf("Имя: %s, Возраст: %d\n", p.name, p.age)
}

// Метод для увеличения возраста на 1 год
func (p *Person) Birthday() {
	p.age++
}

// Структура Circle
type Circle struct {
	radius float64
}

// Метод для вычисления площади круга
func (c Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}

// Интерфейс Shape
type Shape interface {
	Area() float64
}

// Структура Rectangle
type Rectangle struct {
	width, height float64
}

// Реализация метода Area для Rectangle
func (r Rectangle) Area() float64 {
	return r.width * r.height
}

// Функция для вывода площади объектов
func PrintAreas(shapes []Shape) {
	for _, shape := range shapes {
		fmt.Printf("Площадь: %.2f\n", shape.Area())
	}
}

// Интерфейс Stringer
type Stringer interface {
	String() string
}

// Структура Book
type Book struct {
	title  string
	author string
}

// Реализация интерфейса Stringer для Book
func (b Book) String() string {
	return fmt.Sprintf("Книга: %s, автор: %s", b.title, b.author)
}

func main() {
	// Пример работы с Person
	person := Person{name: "Максим", age: 20}
	person.Info()
	person.Birthday()
	person.Info()

	// Пример работы с Circle и Rectangle
	circle := Circle{radius: 5}
	rectangle := Rectangle{width: 8, height: 7}

	// Вывод площади объектов
	shapes := []Shape{circle, rectangle}
	PrintAreas(shapes)

	// Пример работы с Book
	book := Book{title: "The Go Programming Language", author: "Alan A. A. Donovan"}
	fmt.Println(book.String())
}
