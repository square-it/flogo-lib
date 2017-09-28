package property

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRegisterEmptyId register with empty id
func TestRegisterEmptyId(t *testing.T) {
	err := Register("", "a value")
	assert.NotNil(t, err)
}

// TestRegisterOk register ok
func TestRegisterOk(t *testing.T) {
	err := Register("id_test", "a value")
	assert.Nil(t, err)

	value, _ := Get("id_test")
	assert.Equal(t, "a value", value)
}

// TestRegisterDuplicate register duplicate
func TestRegisterDuplicate(t *testing.T) {
	err := Register("id_test2", "a value")
	assert.Nil(t, err)

	err = Register("id_test2", "a value")
	assert.NotNil(t, err)
}

// TestDefaultResolverOk resolves environment and property value
func TestDefaultResolverOk(t *testing.T) {

	Register("id_test", "a value")

	value, _ := Resolve("id_test")
	assert.Equal(t, "a value", value)
}

func TestPropertyResolution(t *testing.T) {
	Register("id_test", "a value")

	value, set := ResolveProperty(nil, "id_test")
	assert.True(t, set)
	assert.Equal(t, "a value", value)
}

func TestPropertyResolutionNotSet(t *testing.T) {
	Register("id_test", "a value")

	_, set := ResolveEnv(nil, "id_test_NS")
	assert.False(t, set)
}

func TestEnvResolution(t *testing.T) {
	os.Setenv("TEST_FLOGO2", "my_test_value2")
	defer os.Unsetenv("TEST_FLOGO2")

	value, set := ResolveEnv(nil, "TEST_FLOGO2")
	assert.True(t, set)
	assert.Equal(t, "my_test_value2", value)
}

func TestEnvResolutionNotSet(t *testing.T) {
	os.Setenv("TEST_FLOGO2", "my_test_value2")
	defer os.Unsetenv("TEST_FLOGO2")

	_, set := ResolveEnv(nil, "TEST_FLOGO_NS")
	assert.False(t, set)
}