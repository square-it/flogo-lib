package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"os"
)

// TestResolveEnvEmpty test an empty property resolution
func TestResolveEnvEmpty(t *testing.T) {
	value := ResolveEnv("")
	assert.Equal(t, "", value)
}

// TestResolveEnvNoEnvironment test an non environment property resolution
func TestResolveEnvNoEnvironment(t *testing.T) {
	value := ResolveEnv("testname")
	assert.Equal(t, "testname", value)

	value = ResolveEnv("a")
	assert.Equal(t, "a", value)

	value = ResolveEnv("{testname")
	assert.Equal(t, "{testname", value)

	value = ResolveEnv("testname}")
	assert.Equal(t, "testname}", value)
}

// TestResolveEnvOk test an environment property resolution
func TestResolveEnvOk(t *testing.T) {
	os.Setenv("TEST_FLOGO", "my_test_value")
	defer os.Unsetenv("TEST_FLOGO")

	value := ResolveEnv("TEST_FLOGO")
	assert.Equal(t, "TEST_FLOGO", value)

	value = ResolveEnv("{TEST_FLOGO}")
	assert.Equal(t, "my_test_value", value)

	value = ResolveEnv("{}")
	assert.Equal(t, "", value)
}