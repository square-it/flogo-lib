package data

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetResolutionDetails(t *testing.T) {

	a := "property.Prop1"
	details, err := GetResolutionDetails(a)
	assert.Nil(t, err)
	assert.Equal(t, "property", details.ResolverName)
	assert.Equal(t, "Prop1", details.Property)
	assert.Equal(t, "", details.Item)
	assert.Equal(t, "", details.Path)

	a = "env.VAR1"
	details, err = GetResolutionDetails(a)
	assert.Nil(t, err)
	assert.Equal(t, "env", details.ResolverName)
	assert.Equal(t, "VAR1", details.Property)
	assert.Equal(t, "", details.Item)
	assert.Equal(t, "", details.Path)

	// Resolution of first level Activity expression
	a = "activity[myactivityId].myAttributeName"
	details, err = GetResolutionDetails(a)
	assert.Nil(t, err)
	assert.Equal(t, "activity", details.ResolverName)
	assert.Equal(t, "myAttributeName", details.Property)
	assert.Equal(t, "myactivityId", details.Item)
	assert.Equal(t, "", details.Path)

	// Resolution of second level Activity expression map
	a = "activity[myactivityId].myMapAttributeName.mapkey"
	details, err = GetResolutionDetails(a)
	assert.Nil(t, err)
	assert.Equal(t, "activity", details.ResolverName)
	assert.Equal(t, "myMapAttributeName", details.Property)
	assert.Equal(t, "myactivityId", details.Item)
	assert.Equal(t, ".mapkey", details.Path)

	// Resolution of second level Activity expression array
	a = "activity.myactivityId.myArrayAttributeName[0]"
	details, err = GetResolutionDetails(a)
	assert.Nil(t, err)
	assert.Equal(t, "activity", details.ResolverName)
	assert.Equal(t, "myArrayAttributeName", details.Property)
	assert.Equal(t, "myactivityId", details.Item)
	assert.Equal(t, "[0]", details.Path)

	// Resolution of first level Activity expression
	a = "activity.myactivityId.myAttributeName"
	details, err = GetResolutionDetails(a)
	assert.Nil(t, err)
	assert.Equal(t, "activity", details.ResolverName)
	assert.Equal(t, "myAttributeName", details.Property)
	assert.Equal(t, "myactivityId", details.Item)
	assert.Equal(t, "", details.Path)

	// Resolution of second level Activity expression map
	a = "activity.myactivityId.myMapAttributeName.mapkey"
	details, err = GetResolutionDetails(a)
	assert.Nil(t, err)
	assert.Equal(t, "activity", details.ResolverName)
	assert.Equal(t, "myMapAttributeName", details.Property)
	assert.Equal(t, "myactivityId", details.Item)
	assert.Equal(t, ".mapkey", details.Path)

	// Resolution of second level Activity expression array
	a = "activity.myactivityId.myArrayAttributeName[0]"
	details, err = GetResolutionDetails(a)
	assert.Nil(t, err)
	assert.Equal(t, "activity", details.ResolverName)
	assert.Equal(t, "myArrayAttributeName", details.Property)
	assert.Equal(t, "myactivityId", details.Item)
	assert.Equal(t, "[0]", details.Path)
}


func TestGetResolutionDetailsOld(t *testing.T) {

	a := "${property.Prop1}"
	details, err := GetResolutionDetailsOld(a)
	assert.Nil(t, err)
	assert.Equal(t, "property", details.ResolverName)
	assert.Equal(t, "Prop1", details.Property)
	assert.Equal(t, "", details.Item)
	assert.Equal(t, "", details.Path)

	a = "${env.VAR1}"
	details, err = GetResolutionDetailsOld(a)
	assert.Nil(t, err)
	assert.Equal(t, "env", details.ResolverName)
	assert.Equal(t, "VAR1", details.Property)
	assert.Equal(t, "", details.Item)
	assert.Equal(t, "", details.Path)

	a = "${trigger.val}"
	details, err = GetResolutionDetailsOld(a)
	assert.Nil(t, err)
	assert.Equal(t, "trigger", details.ResolverName)
	assert.Equal(t, "val", details.Property)
	assert.Equal(t, "", details.Item)
	assert.Equal(t, "", details.Path)

	a = "${trigger.val}.value"
	details, err = GetResolutionDetailsOld(a)
	assert.Nil(t, err)
	assert.Equal(t, "trigger", details.ResolverName)
	assert.Equal(t, "val", details.Property)
	assert.Equal(t, "", details.Item)
	assert.Equal(t, ".value", details.Path)

	a = "${trigger.val}[0]"
	details, err = GetResolutionDetailsOld(a)
	assert.Nil(t, err)
	assert.Equal(t, "trigger", details.ResolverName)
	assert.Equal(t, "val", details.Property)
	assert.Equal(t, "", details.Item)
	assert.Equal(t, "[0]", details.Path)

	// Resolution of first level Activity expression
	a = "${activity.myactivityId.myAttributeName}"
	details, err = GetResolutionDetailsOld(a)
	assert.Nil(t, err)
	assert.Equal(t, "activity", details.ResolverName)
	assert.Equal(t, "myAttributeName", details.Property)
	assert.Equal(t, "myactivityId", details.Item)
	assert.Equal(t, "", details.Path)

	// Resolution of second level Activity expression map
	a = "${activity.myactivityId.myMapAttributeName}.mapkey"
	details, err = GetResolutionDetailsOld(a)
	assert.Nil(t, err)
	assert.Equal(t, "activity", details.ResolverName)
	assert.Equal(t, "myMapAttributeName", details.Property)
	assert.Equal(t, "myactivityId", details.Item)
	assert.Equal(t, ".mapkey", details.Path)

	// Resolution of second level Activity expression array
	a = "${activity.myactivityId.myArrayAttributeName}[0]"
	details, err = GetResolutionDetailsOld(a)
	assert.Nil(t, err)
	assert.Equal(t, "activity", details.ResolverName)
	assert.Equal(t, "myArrayAttributeName", details.Property)
	assert.Equal(t, "myactivityId", details.Item)
	assert.Equal(t, "[0]", details.Path)
}
