package activity

import (
	"testing"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/stretchr/testify/assert"
)

func TestResolver(t *testing.T){
	// Resolve myStringAttribute
	myStringAttribute := data.NewAttribute("${activity.myActivity.myStringAttribute}", data.STRING, "attr value")
	scope := data.NewSimpleScope([]*data.Attribute{myStringAttribute}, nil)

	resolver := newResolver(scope)
	result, ok := resolver.Resolve("${activity.myActivity.myStringAttribute}")
	assert.True(t, ok)
	assert.NotNil(t, result)
	act, ok := result.(string)
	assert.Equal(t,"attr value", act)
	assert.True(t,ok)

}
