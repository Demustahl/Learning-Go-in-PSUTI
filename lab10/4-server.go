package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
)

func main() {
	// Загрузка сертификата сервера
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatalf("Не удалось загрузить сертификат сервера: %v", err)
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
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()

	// Создание TLS-листенера
	listener, err := tls.Listen("tcp", ":8443", tlsConfig)
	if err != nil {
		log.Fatalf("Не удалось создать TLS-листенер: %v", err)
	}
	defer listener.Close()
	fmt.Println("TLS-сервер запущен на порту 8443")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Ошибка при принятии соединения: %v", err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	// Получаем информацию о клиентском сертификате
	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		log.Println("Не удалось привести соединение к tls.Conn")
		return
	}
	err := tlsConn.Handshake()
	if err != nil {
		log.Printf("Ошибка TLS рукопожатия: %v", err)
		return
	}

	state := tlsConn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		log.Println("Нет клиентского сертификата")
		return
	}

	clientCert := state.PeerCertificates[0]
	fmt.Printf("Аутентифицирован клиент: %s\n", clientCert.Subject.CommonName)

	// Чтение данных от клиента
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		log.Printf("Ошибка при чтении данных: %v", err)
		return
	}
	message := string(buf[:n])
	fmt.Printf("Получено сообщение от клиента: %s\n", message)

	// Ответ клиенту
	response := "Сообщение получено сервером"
	conn.Write([]byte(response))
}
