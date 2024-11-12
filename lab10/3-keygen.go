// keygen.go
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

func generateKeys(bits int) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}

	// Сохраняем приватный ключ
	privateKeyFile, err := os.Create("client_private_key.pem")
	if err != nil {
		return err
	}
	defer privateKeyFile.Close()

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return err
	}

	// Сохраняем публичный ключ
	publicKey := &privateKey.PublicKey
	publicKeyFile, err := os.Create("client_public_key.pem")
	if err != nil {
		return err
	}
	defer publicKeyFile.Close()

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}

	publicKeyPEM := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	if err := pem.Encode(publicKeyFile, publicKeyPEM); err != nil {
		return err
	}

	fmt.Println("Ключи сгенерированы и сохранены в файлы 'client_private_key.pem' и 'client_public_key.pem'.")
	return nil
}

func main() {
	err := generateKeys(2048)
	if err != nil {
		log.Fatalf("Ошибка при генерации ключей: %v", err)
	}
}
