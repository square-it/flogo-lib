package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetType(t *testing.T) {

	resType, err := GetType("res://flow:myflow")
	assert.Nil(t, err)
	assert.Equal(t, "flow", resType)
}
