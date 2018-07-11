package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"os"
	"github.com/TIBCOSoftware/flogo-lib/config"
)

func TestSecretKeyDefault(t *testing.T) {
	defer func() {
		SetSecretValueDecoder(nil)
	}()
	decoder := GetSecretValueDecoder()
	decrypted, err := decoder.DecodeValue("QReglAoZB81p3jj9dCXA2oEbauCXu/N8G5wwqQjeARMG")
	assert.Nil(t, err)
	assert.Equal(t, "mysecurepassword3", decrypted)
}

func TestSecretKeyEnv(t *testing.T) {
	os.Setenv(config.ENV_DATA_SECRET_KEY_KEY, "mysecretkey1")
	defer func() {
		os.Unsetenv(config.ENV_DATA_SECRET_KEY_KEY)
		SetSecretValueDecoder(nil)
	}()
	decrypted, err := GetSecretValueDecoder().DecodeValue("Dpn7oxbWZPHlLkPkzgm+qZFHfGHlAKoFXcu5RhbNlQZS")
	assert.Nil(t, err)
	assert.Equal(t, "mysecurepassword1", decrypted)
}

func TestSecretKey(t *testing.T) {
	defer func() {
		SetSecretValueDecoder(nil)
	}()
	SetSecretValueDecoder(&KeyBasedSecretValueDecoder{Key: "mysecretkey2"})
	decrypted, err := GetSecretValueDecoder().DecodeValue("4hmxSs2mjKfExP0lLLW9R9hK8631ce2r5RWDKYJjDZV2")
	assert.Nil(t, err)
	assert.Equal(t, "mysecurepassword2", decrypted)
}



