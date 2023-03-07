package encode

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"strconv"
)

// deshifu returns user id by string from cookies
func Deshifu(a string) (int, error) {
	key := sha256.Sum256([]byte(password))

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return 0, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return 0, err
	}

	nonce := []byte(tinder)

	encrypted, err := hex.DecodeString(a)
	if err != nil {
		return 0, err
	}
	// расшифровываем
	decrypted, err := aesgcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return 0, err
	}
	//fmt.Println(string(decrypted))
	int1, err := strconv.Atoi(string(decrypted))
	if err != nil {
		log.Println("Can not convert this []byte to int")
		return 0, err
	}
	return int1, nil //string(decrypted), nil
}
