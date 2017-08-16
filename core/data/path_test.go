package data

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAttrPath(t *testing.T) {

	a := "sensorData.temp"
	GetAttrPath(a)
	name, path, pt := GetAttrPath(a)
	fmt.Printf("Name: %s, Path: %s, PathType: %d\n", name, path, pt)

	a = "T.v"
	name, path, pt = GetAttrPath(a)
	fmt.Printf("Name: %s, Path: %s, PathType: %d\n", name, path, pt)

	a = "{T1.v}.myAttr"
	name, path, pt = GetAttrPath(a)
	fmt.Printf("Name: %s, Path: %s, PathType: %d\n", name, path, pt)

	a = "{T1.v}[0]"
	name, path, pt = GetAttrPath(a)
	fmt.Printf("Name: %s, Path: %s, PathType: %d\n", name, path, pt)

	a = "v[0]"
	name, path, pt = GetAttrPath(a)
	fmt.Printf("Name: %s, Path: %s, PathType: %d\n", name, path, pt)

	a = "v"
	name, path, pt = GetAttrPath(a)
	fmt.Printf("Name: %s, Path: %s, PathType: %d\n", name, path, pt)

	a = "{T1.v}"
	name, path, pt = GetAttrPath(a)
	fmt.Printf("Name: %s, Path: %s, PathType: %d\n", name, path, pt)


}

func TestGetResolverType(t *testing.T) {
	// Resolution of Property expression
	a := "${property.Prop1}"
	resType, err := GetResolverType(a)
	assert.Nil(t, err)
	assert.Equal(t, RES_PROPERTY, resType)

	// Resolution of Environment expression
	a = "${env.VAR1}"
	resType, err = GetResolverType(a)
	assert.Nil(t, err)
	assert.Equal(t, RES_PROPERTY, resType)

	// Resolution of first level Activity expression
	a = "${activity.myStringAttribute}"
	resType, err = GetResolverType(a)
	assert.Nil(t, err)
	assert.Equal(t, RES_ACTIVITY, resType)

	// Resolution of second level Activity expression
	a = "${activity.myMapAttribute.myMapKey}"
	resType, err = GetResolverType(a)
	assert.Nil(t, err)
	assert.Equal(t, RES_ACTIVITY, resType)

	// Resolution of flat Trigger expression
	a = "${trigger.myTrigger}"
	resType, err = GetResolverType(a)
	assert.Nil(t, err)
	assert.Equal(t, RES_TRIGGER, resType)
}
