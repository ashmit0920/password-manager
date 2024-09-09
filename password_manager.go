package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"io"
	"os"
)

const encryptionKey = "mysecretkey12345"

type PasswordEntry struct {
	Website  string `json:"website"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var passwords []PasswordEntry

func savePassword(entry PasswordEntry, filename string) error {
	passwords = append(passwords, entry)
	data, err := json.Marshal(passwords)
	if err != nil {
		return err
	}

	encryptedData, err := encrypt(data, []byte(encryptionKey))
	if err != nil {
		return err
	}

	return os.WriteFile(filename, encryptedData, 0644)
}

func loadPasswords(filename string) ([]PasswordEntry, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	decryptedData, err := decrypt(data, []byte(encryptionKey))
	if err != nil {
		return nil, err
	}

	var entries []PasswordEntry
	err = json.Unmarshal(decryptedData, &entries)
	return entries, err
}

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
