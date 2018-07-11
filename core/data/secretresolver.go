package data

import (
	"github.com/TIBCOSoftware/flogo-lib/config"
	"crypto/sha256"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
)

var secretDecoder SecretValueDecoder

type SecretValueDecoder interface {
	DecodeValue(value interface{}) (string, error)
}


func SetSecretValueDecoder(pwdResolver SecretValueDecoder ) {
	secretDecoder = pwdResolver
}

func GetSecretValueDecoder() SecretValueDecoder {
	if secretDecoder == nil {
		secretDecoder = &KeyBasedSecretValueDecoder{Key: config.GetSecretKey()}
	}
	return secretDecoder
}


type KeyBasedSecretValueDecoder struct {
	Key string
}

func (defaultResolver *KeyBasedSecretValueDecoder) DecodeValue(value interface{}) (string, error) {
	if value != nil   {
		fmt.Println(defaultResolver.Key)
		if defaultResolver.Key != "" {
			keyBytes := sha256.Sum256([]byte(defaultResolver.Key))
			return decrypt(keyBytes[:], value.(string))
		}
		return value.(string), nil
	}
	return "", nil
}


// decrypt from base64 to decrypted string
func decrypt(key []byte, cryptoText string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return fmt.Sprintf("%s", ciphertext), nil
}


