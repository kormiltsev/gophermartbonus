package encode

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

// shifu returns cookie string by user id
func Shifu(a int) (string, error) {
	key := sha256.Sum256([]byte(password))

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := []byte(tinder)

	ciphertext := aesgcm.Seal(nil, nonce, []byte(strconv.Itoa(a)), nil)

	export := hex.EncodeToString(ciphertext)
	return export, nil
}
