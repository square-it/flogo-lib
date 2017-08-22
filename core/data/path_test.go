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


	// Resolution of old Trigger activity expression
	a = "{T.pathParams}.myParam"
	name, path, pt = GetAttrPath(a)
	assert.Equal(t, "{T.pathParams}", name)
	assert.Equal(t, "myParam", path)
	assert.Equal(t, PT_MAP, pt)

	// Resolution of activity expression
	a = "${activity.my_activityid.mystring}"
	name, path, pt = GetAttrPath(a)
	assert.Equal(t, "${activity.my_activityid.mystring}", name)
	assert.Equal(t, "", path)
	assert.Equal(t, PT_SIMPLE, pt)

	// Resolution of activity expression
	a = "${activity.my_activityid.mymap}.mypath"
	name, path, pt = GetAttrPath(a)
	assert.Equal(t, "${activity.my_activityid.mymap}", name)
	assert.Equal(t, "mypath", path)
	assert.Equal(t, PT_MAP, pt)

	// Resolution of activity expression
	a = "${activity.my_activityid.myarray}[0]"
	name, path, pt = GetAttrPath(a)
	assert.Equal(t, "${activity.my_activityid.myarray}", name)
	assert.Equal(t, "0", path)
	assert.Equal(t, PT_ARRAY, pt)

}

func TestGetResolverType(t *testing.T) {
	// Resolution of Old Trigger expression
	a := "{T.pathParams}.myParam"
	resType, err := GetResolverType(a)
	assert.Nil(t, err)
	assert.Equal(t, RES_DEFAULT, resType)

	// Resolution of Property expression
	a = "${property.Prop1}"
	resType, err = GetResolverType(a)
	assert.Nil(t, err)
	assert.Equal(t, RES_PROPERTY, resType)

	// Resolution of Environment expression
	a = "${env.VAR1}"
	resType, err = GetResolverType(a)
	assert.Nil(t, err)
	assert.Equal(t, RES_PROPERTY, resType)

	// Resolution of first level Activity expression
	a = "${activity.myactivityId.myAttributeName}"
	resType, err = GetResolverType(a)
	assert.Nil(t, err)
	assert.Equal(t, RES_ACTIVITY, resType)

	// Resolution of second level Activity expression map
	a = "${activity.myactivityId.myMapAttributeName}.mapkey"
	resType, err = GetResolverType(a)
	assert.Nil(t, err)
	assert.Equal(t, RES_ACTIVITY, resType)

	// Resolution of second level Activity expression array
	a = "${activity.myactivityId.myMapAttributeName}[0]"
	resType, err = GetResolverType(a)
	assert.Nil(t, err)
	assert.Equal(t, RES_ACTIVITY, resType)

	// Resolution of flat Trigger expression
	a = "${trigger.myTrigger}"
	resType, err = GetResolverType(a)
	assert.Nil(t, err)
	assert.Equal(t, RES_TRIGGER, resType)
}
