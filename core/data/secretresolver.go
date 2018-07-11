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

// SecretValueDecoder defines method for decoding value
type SecretValueDecoder interface {
	DecodeValue(value interface{}) (string, error)
}

// Set secret value decoder
func SetSecretValueDecoder(pwdResolver SecretValueDecoder ) {
	secretDecoder = pwdResolver
}

// Get secret value decoder. If not already set by SetSecretValueDecoder(), will return default KeyBasedSecretValueDecoder
// where decoding key value is expected to be set through FLOGO_DATA_SECRET_KEY environment variable.
// If key is not set, a default key value(github.com/TIBCOSoftware/flogo-lib/config.ENV_DATA_SECRET_KEY_DEFAULT) will be used.
func GetSecretValueDecoder() SecretValueDecoder {
	if secretDecoder == nil {
		secretDecoder = &KeyBasedSecretValueDecoder{Key: config.GetDataSecretKey()}
	}
	return secretDecoder
}

// A key based secret value decoder. Secret value encryption/decryption is based on SHA256
// and refers https://gist.github.com/willshiao/f4b03650e5a82561a460b4a15789cfa1
type KeyBasedSecretValueDecoder struct {
	Key string
}

// Decode value based on a key
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


