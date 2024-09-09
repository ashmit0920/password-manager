package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type PasswordEntry struct {
	Website  string `json:"website"`
	Username string `json:"username"`
	Password string `json:"password"`
}

const key = "mysecretkey12345" // Must be 16, 24, or 32 bytes for AES-128, 192, or 256

// Encrypts text with AES encryption
func encrypt(text []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], text)

	return ciphertext, nil
}

// Decrypts AES encrypted text
func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

// Save the password entry
func savePassword(entry PasswordEntry, filename string) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	encryptedData, err := encrypt(data, []byte(key))
	if err != nil {
		return err
	}

	return os.WriteFile(filename, encryptedData, 0644)
}

// Retrieve and decrypt the password entry
func retrievePassword(filename string) (PasswordEntry, error) {
	var entry PasswordEntry
	data, err := os.ReadFile(filename)
	if err != nil {
		return entry, err
	}

	decryptedData, err := decrypt(data, []byte(key))
	if err != nil {
		return entry, err
	}

	err = json.Unmarshal(decryptedData, &entry)
	return entry, err
}

func main() {
	var entry PasswordEntry
	fmt.Println("Enter website: ")
	fmt.Scanln(&entry.Website)
	fmt.Println("Enter username: ")
	fmt.Scanln(&entry.Username)
	fmt.Println("Enter password: ")
	fmt.Scanln(&entry.Password)

	filename := "passwords.enc"
	if err := savePassword(entry, filename); err != nil {
		fmt.Println("Error saving password:", err)
	} else {
		fmt.Println("Password saved successfully!")
	}

	// Retrieve and display password
	savedEntry, err := retrievePassword(filename)
	if err != nil {
		fmt.Println("Error retrieving password:", err)
	} else {
		fmt.Printf("Saved Entry: %v\n", savedEntry)
	}
}
