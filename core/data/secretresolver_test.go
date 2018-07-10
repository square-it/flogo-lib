package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"os"
	"github.com/TIBCOSoftware/flogo-lib/config"
)

func TestSecretKeyEnv(t *testing.T) {
	os.Setenv(config.ENV_SECRET_KEY_KEY, "mysecretkey1")
	defer os.Unsetenv(config.ENV_SECRET_KEY_KEY)

	decoder := GetSecretValueDecoder()
	decrypted, err := decoder.DecodeValue("Dpn7oxbWZPHlLkPkzgm+qZFHfGHlAKoFXcu5RhbNlQZS")
	assert.Nil(t, err)
	assert.Equal(t, "mysecurepassword1", decrypted)
}

func TestSecretKey(t *testing.T) {
	decoder := &KeyBasedSecretValueDecoder{Key: "mysecretkey2"}
	SetSecretValueDecoder(decoder)
	decrypted, err := decoder.DecodeValue("4hmxSs2mjKfExP0lLLW9R9hK8631ce2r5RWDKYJjDZV2")
	assert.Nil(t, err)
	assert.Equal(t, "mysecurepassword2", decrypted)
}

