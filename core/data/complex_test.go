package data

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshComplexObject(t *testing.T) {
	complexStr := `{"metadata":"this is metdata string","value":""}`
	complexObject, err := CoerceToComplexObject(complexStr)
	assert.Nil(t, err)
	assert.NotNil(t, complexObject)
	assert.NotEqual(t, "", complexObject.Value)
}

