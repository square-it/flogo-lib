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

	// Resolution of Property expression
	a = "${property.Prop1}"
	name, path, pt = GetAttrPath(a)
	fmt.Printf("Name: %s, Path: %s, PathType: %d\n", name, path, pt)
	assert.Equal(t, "property", name)
	assert.Equal(t, "Prop1", path)
	assert.Equal(t, pt, PT_PROPERTY)

	// Resolution of Environment expression
	a = "${env.VAR1}"
	name, path, pt = GetAttrPath(a)
	fmt.Printf("Name: %s, Path: %s, PathType: %d\n", name, path, pt)
	assert.Equal(t, "env", name)
	assert.Equal(t, "VAR1", path)
	assert.Equal(t, pt, PT_PROPERTY)

	// Resolution of flat Activity expression
	a = "${activity.myActivity}"
	name, path, pt = GetAttrPath(a)
	fmt.Printf("Name: %s, Path: %s, PathType: %d\n", name, path, pt)
	assert.Equal(t, "activity", name)
	assert.Equal(t, "myActivity", path)
	assert.Equal(t, pt, PT_ACTIVITY)

	// Resolution of flat Trigger expression
	a = "${trigger.myTrigger}"
	name, path, pt = GetAttrPath(a)
	fmt.Printf("Name: %s, Path: %s, PathType: %d\n", name, path, pt)
	assert.Equal(t, "trigger", name)
	assert.Equal(t, "myTrigger", path)
	assert.Equal(t, pt, PT_TRIGGER)
}
