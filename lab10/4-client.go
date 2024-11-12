package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	// Загрузка сертификата клиента
	cert, err := tls.LoadX509KeyPair("client.crt", "client.key")
	if err != nil {
		log.Fatalf("Не удалось загрузить сертификат клиента: %v", err)
	}

	// Загрузка корневого сертификата (CA)
	caCert, err := ioutil.ReadFile("ca.crt")
	if err != nil {
		log.Fatalf("Не удалось загрузить сертификат CA: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Настройка TLS-конфигурации
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: false,
	}
	tlsConfig.BuildNameToCertificate()

	// Установка соединения с сервером
	conn, err := tls.Dial("tcp", "localhost:8443", tlsConfig)
	if err != nil {
		log.Fatalf("Не удалось установить TLS-соединение: %v", err)
	}
	defer conn.Close()

	fmt.Println("TLS-соединение с сервером установлено")

	// Отправка данных серверу
	message := "Привет от клиента"
	_, err = conn.Write([]byte(message))
	if err != nil {
		log.Fatalf("Ошибка при отправке данных: %v", err)
	}

	// Чтение ответа от сервера
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalf("Ошибка при чтении данных: %v", err)
	}
	response := string(buf[:n])
	fmt.Printf("Ответ от сервера: %s\n", response)
}
