package activity

import (
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolver(t *testing.T) {
	// Resolve myStringAttribute
	myStringAttribute := data.NewAttribute("_A.myActivity.myStringAttribute", data.STRING, "attr value")
	scope := data.NewSimpleScope([]*data.Attribute{myStringAttribute}, nil)

	resolver := data.GetResolver(data.RES_ACTIVITY)
	result, ok := resolver(scope, "myActivity.myStringAttribute")
	assert.True(t, ok)
	assert.NotNil(t, result)
	act, ok := result.(string)
	assert.Equal(t, "attr value", act)
	assert.True(t, ok)
}

func TestResolution(t *testing.T) {
	// Resolve myStringAttribute
	myStringAttribute := data.NewAttribute("_A.myActivity.myStringAttribute", data.STRING, "attr value")
	scope := data.NewSimpleScope([]*data.Attribute{myStringAttribute}, nil)

	resType, attrName, path, err  := data.GetResolutionInfo("${activity.myActivity.myStringAttribute}")
	assert.Nil(t, err)
	assert.Equal(t, data.RES_ACTIVITY, resType)
	assert.Equal(t, "", path)

	resolver := data.GetResolver(resType)
	result, ok := resolver(scope, attrName)
	assert.True(t, ok)
	assert.NotNil(t, result)
	act, ok := result.(string)
	assert.Equal(t, "attr value", act)
	assert.True(t, ok)
}



