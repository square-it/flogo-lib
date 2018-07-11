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
		if defaultResolver.Key != "" {
			kBytes := sha256.Sum256([]byte(defaultResolver.Key))
			return decryptValue(kBytes[:], value.(string))
		}
		return value.(string), nil
	}
	return "", nil
}


// decrypt from base64 to decrypted string
func decryptValue(key []byte, encryptedData string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
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


