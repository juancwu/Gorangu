package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	"golang.org/x/crypto/scrypt"
)

func main() {
    var (
        password = []byte("secret password")
        data = []byte("need to protect this data")
    )

    ciphertext, err := encrypt(password, data)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("ciphertext: %s\n", hex.EncodeToString(ciphertext))

    plaintext, err := decrypt(password, ciphertext)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("plaintext: %s\n", plaintext)
}

func encrypt(password []byte, data []byte) ([]byte, error) {
    key, salt, err := deriveKey(password, nil)
    if err != nil {
        return nil, err
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err = rand.Read(nonce); err != nil {
        return nil, err
    }

    cipherText := gcm.Seal(nonce, nonce, data, nil)
    cipherText = append(cipherText, salt...)

    return cipherText, nil
}

func decrypt(password []byte, data []byte) ([]byte, error) {
    salt, data := data[len(data)-32:], data[:len(data)-32]
    key, _, err := deriveKey(password, salt)
    if err != nil {
        return nil, err
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, err
    }

    return plaintext, nil
}

func deriveKey(password []byte, salt []byte) ([]byte, []byte, error) {
    if salt == nil {
        salt = make([]byte, 32)
        if _, err := rand.Read(salt); err != nil {
            return nil, nil, err
        }
    }

    key, err := scrypt.Key(password, salt, 1048576, 8, 1, 32)
    if err != nil {
        return nil, nil, err
    }

    return key, salt, nil
}
