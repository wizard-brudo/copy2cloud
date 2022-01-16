package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// внутренности crypt
var (
	cryptKey = []byte{
		0x9c, 0x93, 0x5b, 0x48, 0x73, 0x0a, 0x55, 0x4d,
		0x6b, 0xfd, 0x7c, 0x63, 0xc8, 0x86, 0xa9, 0x2b,
		0xd3, 0x90, 0x19, 0x8e, 0xb8, 0x12, 0x8a, 0xfb,
		0xf4, 0xde, 0x16, 0x2b, 0x8b, 0x95, 0xf6, 0x38,
	}
	cryptBlock cipher.Block
	cryptRand  = rand.Reader
)

// crypt преобразуется из входа в выход с помощью iv под AES-CTR.
//
// вход и выход могут быть одним и тем же буфером.
//
// Обратите внимание, что шифрование и дешифрование — это одна и та же операция
func crypt(out, in, iv []byte) error {
	if cryptBlock == nil {
		var err error
		cryptBlock, err = aes.NewCipher(cryptKey)
		if err != nil {
			return err
		}
	}
	stream := cipher.NewCTR(cryptBlock, iv)
	stream.XORKeyStream(out, in)
	return nil
}

// Скрыть значение
func Obscure(x string) (string, error) {
	plaintext := []byte(x)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(cryptRand, iv); err != nil {
		return "", errors.New("Не удалось прочитать")
	}
	if err := crypt(ciphertext[aes.BlockSize:], plaintext, iv); err != nil {
		return "", errors.New("Зашифровать не удалось")
	}
	return base64.RawURLEncoding.EncodeToString(ciphertext), nil
}

// MustObscure скрывает значение
func MustObscure(x string) (string, error) {
	out, err := Obscure(x)
	if err != nil {
		return "", err
	}
	return out, nil
}

// Выявить скрытое значение
func Reveal(x string) (string, error) {
	ciphertext, err := base64.RawURLEncoding.DecodeString(x)
	if err != nil {
		return "", errors.New("Декодирование base64 не удалось при раскрытии пароля - он скрыт?")
	}
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("Ввод слишком короткий при раскрытии пароля - он скрыт?")
	}
	buf := ciphertext[aes.BlockSize:]
	iv := ciphertext[:aes.BlockSize]
	if err := crypt(buf, buf, iv); err != nil {
		return "", errors.New("Расшифровка не удалась при раскрытии пароля - он скрыт?")
	}
	return string(buf), nil
}

// MustReveal показывает скрытое значение
func MustReveal(x string) (string, error) {
	out, err := Reveal(x)
	if err != nil {
		return "", err
	}
	return out, nil
}
