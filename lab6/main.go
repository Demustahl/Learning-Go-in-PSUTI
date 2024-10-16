package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

// Главная функция
func main() {
	fmt.Println("Выберите номер задания (1-6):")
	var task int
	fmt.Scan(&task)

	switch task {
	case 1:
		task1()
	case 2:
		task2()
	case 3:
		task3()
	case 4:
		task4()
	case 5:
		task5()
	case 6:
		task6()
	default:
		fmt.Println("Неверный номер задания")
	}
}

// Задание 1: Создание и запуск горутин
func task1() {
	var wg sync.WaitGroup
	wg.Add(3)

	// Запуск функции расчёта факториала
	go func() {
		defer wg.Done()
		calculateFactorial(5)
	}()

	// Запуск функции генерации случайных чисел
	go func() {
		defer wg.Done()
		generateRandomNumbers(5)
	}()

	// Запуск функции вычисления суммы числового ряда
	go func() {
		defer wg.Done()
		calculateSeriesSum(5)
	}()

	// Ожидание завершения всех горутин
	wg.Wait()
}

// Функция для расчёта факториала числа n
func calculateFactorial(n int) {
	result := 1
	for i := 1; i <= n; i++ {
		result *= i
		fmt.Printf("[Факториал] Текущее значение факториала для %d: %d\n", i, result)
		time.Sleep(400 * time.Millisecond) // Имитация задержки
	}
	fmt.Printf("[Факториал] Факториал числа %d равен %d\n", n, result)
}

// Функция для генерации n случайных чисел
func generateRandomNumbers(n int) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		num := rand.Intn(100)
		fmt.Printf("[Случайные числа] Сгенерировано число: %d\n", num)
		time.Sleep(600 * time.Millisecond) // Имитация задержки
	}
	fmt.Println("[Случайные числа] Генерация завершена")
}

// Функция для вычисления суммы числового ряда до n
func calculateSeriesSum(n int) {
	sum := 0
	for i := 1; i <= n; i++ {
		sum += i
		fmt.Printf("[Сумма ряда] Текущее значение суммы: %d\n", sum)
		time.Sleep(500 * time.Millisecond) // Имитация задержки
	}
	fmt.Printf("[Сумма ряда] Сумма числового ряда до %d равна %d\n", n, sum)
}

// Задание 2: Использование каналов для передачи данных
func task2() {
	ch := make(chan int)
	go fibonacciGenerator(10, ch)
	go fibonacciPrinter(ch)

	// Чтобы main не завершился раньше
	time.Sleep(2 * time.Second)
}

// Генератор чисел Фибоначчи
func fibonacciGenerator(n int, ch chan int) {
	defer close(ch) // Закрываем канал после отправки всех чисел
	a, b := 0, 1
	for i := 0; i < n; i++ {
		ch <- a
		a, b = b, a+b
		time.Sleep(100 * time.Millisecond) // Имитация задержки
	}
}

// Печать чисел из канала
func fibonacciPrinter(ch chan int) {
	for num := range ch {
		fmt.Printf("Получено число Фибоначчи: %d\n", num)
	}
	fmt.Println("Канал закрыт, чтение завершено")
}

// Задание 3: Применение select для управления каналами
func task3() {
	numbers := make(chan int)
	messages := make(chan string)

	// Горутин для генерации случайных чисел
	go func() {
		rand.Seed(time.Now().UnixNano())
		for {
			num := rand.Intn(100)
			numbers <- num
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// Горутин для отправки сообщений о чётности/нечётности
	go func() {
		rand.Seed(time.Now().UnixNano())
		for {
			num := rand.Intn(100)
			if num%2 == 0 {
				messages <- fmt.Sprintf("Число %d - чётное", num)
			} else {
				messages <- fmt.Sprintf("Число %d - нечётное", num)
			}
			time.Sleep(700 * time.Millisecond)
		}
	}()

	// Главная горутина использует select для приёма данных из обоих каналов
	go func() {
		for {
			select {
			case num := <-numbers:
				fmt.Printf("Получено число: %d\n", num)
			case msg := <-messages:
				fmt.Println(msg)
			}
		}
	}()

	// Чтобы main не завершился сразу
	time.Sleep(5 * time.Second)
}

// Задание 4: Синхронизация с помощью мьютексов
func task4() {
	counter := 0
	var wg sync.WaitGroup
	var mutex sync.Mutex

	fmt.Println("Включить мьютекс? (y/n):")
	var input string
	fmt.Scan(&input)
	useMutex := input == "y"

	numGoroutines := 100
	incrementsPerGoroutine := 1000

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				if useMutex {
					mutex.Lock()
				}
				counter++
				if useMutex {
					mutex.Unlock()
				}
			}
		}()
	}
	wg.Wait()
	expected := numGoroutines * incrementsPerGoroutine
	fmt.Printf("Ожидаемое значение счётчика: %d\n", expected)
	fmt.Printf("Реальное значение счётчика: %d\n", counter)
	if useMutex {
		fmt.Println("Мьютекс использовался. Гонка данных предотвращена.")
	} else {
		fmt.Println("Мьютекс не использовался. Возможна гонка данных.")
	}
}

// Задание 5: Разработка многопоточного калькулятора
type Request struct {
	a, b   float64
	op     string
	result chan float64
}

func task5() {
	requests := make(chan Request)

	// Запуск серверной части калькулятора
	go calculatorServer(requests)

	// Клиентские запросы
	for i := 0; i < 5; i++ {
		go func(i int) {
			a := float64(i + 1)
			b := float64((i + 1) * 2)
			ops := []string{"+", "-", "*", "/"}
			op := ops[i%4]
			resultChan := make(chan float64)
			req := Request{a, b, op, resultChan}
			requests <- req
			result := <-resultChan
			fmt.Printf("Результат операции %f %s %f = %f\n", a, op, b, result)
		}(i)
	}

	// Чтобы main не завершился раньше
	time.Sleep(2 * time.Second)
}

// Серверная часть калькулятора
func calculatorServer(requests chan Request) {
	for req := range requests {
		var res float64
		switch req.op {
		case "+":
			res = req.a + req.b
		case "-":
			res = req.a - req.b
		case "*":
			res = req.a * req.b
		case "/":
			if req.b != 0 {
				res = req.a / req.b
			} else {
				fmt.Println("Деление на ноль")
				res = 0
			}
		default:
			fmt.Println("Неизвестная операция")
			res = 0
		}
		req.result <- res
	}
}

// Задание 6: Создание пула воркеров с чтением из файла
func task6() {
	// Создаём файл с случайными фразами
	filename := "phrases.txt"
	numPhrases := 20 // Изменено на 20
	err := createRandomPhrasesFile(filename, numPhrases)
	if err != nil {
		fmt.Println("Ошибка при создании файла с фразами:", err)
		return
	}

	// Читаем строки из файла
	lines, err := readLinesFromFile(filename)
	if err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		return
	}

	// Запросить количество воркеров у пользователя
	fmt.Println("Введите количество воркеров:")
	var numWorkers int
	fmt.Scan(&numWorkers)

	tasks := make(chan string, len(lines))
	results := make(chan string, len(lines))
	var wg sync.WaitGroup

	// Запуск воркеров
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, tasks, results, &wg)
	}

	// Запуск горутины для получения и вывода результатов
	go func() {
		for result := range results {
			fmt.Println("Результат:", result)
		}
	}()

	// Отправка задач
	for _, line := range lines {
		tasks <- line
	}
	close(tasks)

	// Ожидание завершения воркеров
	wg.Wait()
	close(results)
}

// Функция воркера
func worker(id int, tasks <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	var i int = 1
	for line := range tasks {
		// Реверсирование строки
		reversed := reverseString(line)
		time.Sleep(500 * time.Millisecond) // Имитация обработки
		results <- fmt.Sprintf("Воркер %d обработал строку: %s", id+1, reversed)
		i++
	}
}

// Функция для реверсирования строки
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Функция для создания файла с случайными фразами
func createRandomPhrasesFile(filename string, numPhrases int) error {
	phrases := []string{
		"Случайная фраза один",
		"Ещё одна случайная фраза",
		"Пример фразы для теста",
		"Генерация фраз для файла",
		"Тестирование программы на Go",
		"Проверка работы воркеров",
		"Обработка строк в горутинах",
		"Параллельное программирование",
		"Синхронизация данных",
		"Каналы и горутины в Go",
		"Многопоточная обработка",
		"Производительность приложений",
		"Оптимизация кода",
		"Структуры данных",
		"Алгоритмы сортировки",
		"Горутины и конкурентность",
		"Протоколы обмена данными",
		"Интерфейсы в Go",
		"Модульное тестирование",
		"Профилирование приложений",
	}

	// Создаём файл
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Записываем случайные фразы в файл
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < numPhrases; i++ {
		index := rand.Intn(len(phrases))
		_, err := file.WriteString(phrases[index] + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// Функция для чтения строк из файла
func readLinesFromFile(filename string) ([]string, error) {
	var lines []string

	// Открываем файл
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Читаем строки из файла
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Проверяем на ошибки сканера
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
