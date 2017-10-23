package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//func TestGetAttrPath(t *testing.T) {
//
//	a := "sensorData.temp"
//	GetAttrPath(a)
//	name, path, pt := GetAttrPath(a)
//	fmt.Printf("Name: %s, Path: %s, PathType: %d\n", name, path, pt)
//
//	a = "T.v"
//	name, path, pt = GetAttrPath(a)
//	fmt.Printf("Name: %s, Path: %s, PathType: %d\n", name, path, pt)
//
//	a = "{T1.v}.myAttr"
//	name, path, pt = GetAttrPath(a)
//	fmt.Printf("Name: %s, Path: %s, PathType: %d\n", name, path, pt)
//
//	a = "{T1.v}[0]"
//	name, path, pt = GetAttrPath(a)
//	fmt.Printf("Name: %s, Path: %s, PathType: %d\n", name, path, pt)
//
//	a = "v[0]"
//	name, path, pt = GetAttrPath(a)
//	fmt.Printf("Name: %s, Path: %s, PathType: %d\n", name, path, pt)
//
//	a = "v"
//	name, path, pt = GetAttrPath(a)
//	fmt.Printf("Name: %s, Path: %s, PathType: %d\n", name, path, pt)
//
//	a = "{T1.v}"
//	name, path, pt = GetAttrPath(a)
//	fmt.Printf("Name: %s, Path: %s, PathType: %d\n", name, path, pt)
//
//	// Resolution of old Trigger activity expression
//	a = "{T.pathParams}.myParam"
//	name, path, pt = GetAttrPath(a)
//	assert.Equal(t, "{T.pathParams}", name)
//	assert.Equal(t, "myParam", path)
//	assert.Equal(t, PT_MAP, pt)
//
//	// Resolution of activity expression
//	a = "${activity.my_activityid.mystring}"
//	name, path, pt = GetAttrPath(a)
//	assert.Equal(t, "${activity.my_activityid.mystring}", name)
//	assert.Equal(t, "", path)
//	assert.Equal(t, PT_SIMPLE, pt)
//
//	// Resolution of activity expression
//	a = "${activity.my_activityid.mymap}.mypath"
//	name, path, pt = GetAttrPath(a)
//	assert.Equal(t, "${activity.my_activityid.mymap}", name)
//	assert.Equal(t, "mypath", path)
//	assert.Equal(t, PT_MAP, pt)
//
//	// Resolution of activity expression
//	a = "${activity.my_activityid.myarray}[0]"
//	name, path, pt = GetAttrPath(a)
//	assert.Equal(t, "${activity.my_activityid.myarray}", name)
//	assert.Equal(t, "0", path)
//	assert.Equal(t, PT_ARRAY, pt)
//
//}

func TestGetResolverType(t *testing.T) {

	// Resolution of Old Trigger expression
	a := "{T.pathParams}.myParam"
	resType, err := GetResolverType(a)
	assert.Nil(t, err)
	assert.Equal(t, RES_SCOPE, resType)

	// Resolution of Property expression
	a = "${property.Prop1}"
	resType, err = GetResolverType(a)
	assert.Nil(t, err)
	assert.Equal(t, RES_PROPERTY, resType)

	// Resolution of Environment expression
	a = "${env.VAR1}"
	resType, err = GetResolverType(a)
	assert.Nil(t, err)
	assert.Equal(t, RES_ENV, resType)

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

	// Resolution of flat Trigger expression
	a = "${myVar}"
	resType, err = GetResolverType(a)
	assert.Nil(t, err)
	assert.Equal(t, RES_SCOPE, resType)

	// Resolution of flat Trigger expression
	a = "myVar"
	resType, err = GetResolverType(a)
	assert.Nil(t, err)
	assert.Equal(t, RES_SCOPE, resType)
}

func TestGetResolutionInfo(t *testing.T) {
	// Resolution of Old Trigger expression
	a := "{T.pathParams}.myParam"
	resType, toResolve, path, err := GetResolutionInfo(a)
	assert.Nil(t, err)
	assert.Equal(t, RES_SCOPE, resType)

	// Resolution of Property expression
	a = "${property.Prop1}"
	resType, toResolve, path, err = GetResolutionInfo(a)
	assert.Nil(t, err)
	assert.Equal(t, RES_PROPERTY, resType)
	assert.Equal(t, "Prop1", toResolve)
	assert.Equal(t, "", path)

	// Resolution of Environment expression
	a = "${env.VAR1}"
	resType, toResolve, path, err = GetResolutionInfo(a)
	assert.Nil(t, err)
	assert.Equal(t, RES_ENV, resType)
	assert.Equal(t, "VAR1", toResolve)
	assert.Equal(t, "", path)

	// Resolution of first level Activity expression
	a = "${activity.myactivityId.myAttributeName}"
	resType, toResolve, path, err = GetResolutionInfo(a)
	assert.Nil(t, err)
	assert.Equal(t, RES_ACTIVITY, resType)
	assert.Equal(t, "myactivityId.myAttributeName", toResolve)
	assert.Equal(t, "", path)

	// Resolution of second level Activity expression map
	a = "${activity.myactivityId.myMapAttributeName}.mapkey"
	resType, toResolve, path, err = GetResolutionInfo(a)
	assert.Nil(t, err)
	assert.Equal(t, RES_ACTIVITY, resType)
	assert.Equal(t, "myactivityId.myMapAttributeName", toResolve)
	assert.Equal(t, ".mapkey", path)

	// Resolution of second level Activity expression array
	a = "${activity.myactivityId.myMapAttributeName}[0]"
	resType, toResolve, path, err = GetResolutionInfo(a)
	assert.Nil(t, err)
	assert.Equal(t, RES_ACTIVITY, resType)
	assert.Equal(t, "myactivityId.myMapAttributeName", toResolve)
	assert.Equal(t, "[0]", path)

	// Resolution of flat Trigger expression
	a = "${trigger.myTrigger}"
	resType, toResolve, path, err = GetResolutionInfo(a)
	assert.Nil(t, err)
	assert.Equal(t, RES_TRIGGER, resType)
	assert.Equal(t, "myTrigger", toResolve)
	assert.Equal(t, "", path)
}

func TestPathGetValue(t *testing.T) {
	// Resolution of Old Trigger expression

	mapVal,_ := CoerceToObject("{\"myParam\":5}")
	path := ".myParam"
	newVal,err := PathGetValue(mapVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 5.0, newVal)

	// Resolution of Old Trigger expression
	arrVal,_ := CoerceToArray("[1,6,3]")
	path = "[1]"
	newVal,err = PathGetValue(arrVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 6.0, newVal)


	mapVal,_ = CoerceToObject("{\"myParam\":{\"nestedMap\":1}}")
	path = ".myParam.nestedMap"
	newVal,err = PathGetValue(mapVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 1.0, newVal)

	arrVal,_ = CoerceToArray("[{\"nestedMap1\":1},{\"nestedMap2\":2}]")
	path = "[1].nestedMap2"
	newVal,err = PathGetValue(arrVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 2.0, newVal)

	mapVal,_ = CoerceToObject("{\"myParam\":{\"nestedArray\":[7,8,9]}}")
	path = ".myParam.nestedArray[1]"
	newVal,err = PathGetValue(mapVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 8.0, newVal)

	arrVal,_ = CoerceToArray("[{\"nestedMap1\":1},{\"nestedMap2\":{\"nestedArray\":[7,8,9]}}]")
	path = "[1].nestedMap2.nestedArray[2]"
	newVal,err = PathGetValue(arrVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 9.0, newVal)

	mapVal,_ = CoerceToObject("{\"myParam\":{\"nestedArray\":[7,8,9]}}")
	path = ".myParam.nestedArray"
	newVal,err = PathGetValue(mapVal, path)
	assert.Nil(t, err)
	//todo check if array

	arrVal,_ = CoerceToArray("[{\"nestedMap1\":1},{\"nestedMap2\":{\"nestedArray\":[7,8,9]}}]")
	path = "[1].nestedMap2"
	newVal,err = PathGetValue(arrVal, path)
	assert.Nil(t, err)
	//todo check if map

}

func TestPathSetValue(t *testing.T) {
	// Resolution of Old Trigger expression

	mapVal,_ := CoerceToObject("{\"myParam\":5}")
	path := ".myParam"
	err := PathSetValue(mapVal, path, 6)
	assert.Nil(t, err)
	newVal,err := PathGetValue(mapVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 6, newVal)

	// Resolution of Old Trigger expression
	arrVal,_ := CoerceToArray("[1,6,3]")
	path = "[1]"
	err = PathSetValue(arrVal, path, 4)
	assert.Nil(t, err)
	newVal,err = PathGetValue(arrVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 4, newVal)

	mapVal,_ = CoerceToObject("{\"myParam\":{\"nestedMap\":1}}")
	path = ".myParam.nestedMap"
	assert.Nil(t, err)
	err = PathSetValue(mapVal, path, 7)
	newVal,err = PathGetValue(mapVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 7, newVal)

	arrVal,_ = CoerceToArray("[{\"nestedMap1\":1},{\"nestedMap2\":2}]")
	path = "[1].nestedMap2"
	err = PathSetValue(arrVal, path, 3)
	assert.Nil(t, err)
	newVal,err = PathGetValue(arrVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 3, newVal)


	mapVal,_ = CoerceToObject("{\"myParam\":{\"nestedArray\":[7,8,9]}}")
	path = ".myParam.nestedArray[1]"
	err = PathSetValue(mapVal, path, 1)
	assert.Nil(t, err)
	newVal,err = PathGetValue(mapVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 1, newVal)

	arrVal,_ = CoerceToArray("[{\"nestedMap1\":1},{\"nestedMap2\":{\"nestedArray\":[7,8,9]}}]")
	path = "[1].nestedMap2.nestedArray[2]"
	err = PathSetValue(arrVal, path, 5)
	assert.Nil(t, err)
	newVal,err = PathGetValue(arrVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 5, newVal)

	//mapVal,_ = CoerceToObject("{\"myParam\":{\"nestedArray\":[7,8,9]}}")
	//path = ".myParam.nestedArray"
	//err = PathSetValue(arrVal, path, 3)
	//assert.Nil(t, err)
	//newVal,err = PathGetValue(mapVal, path)
	//assert.Nil(t, err)
	////todo check if array
	//
	//arrVal,_ = CoerceToArray("[{\"nestedMap1\":1},{\"nestedMap2\":{\"nestedArray\":[7,8,9]}}]")
	//path = "[1].nestedMap2"
	//assert.Nil(t, err)
	//err = PathSetValue(arrVal, path, 3)
	//newVal,err = PathGetValue(arrVal, path)
	//assert.Nil(t, err)
	//////todo check if map
}