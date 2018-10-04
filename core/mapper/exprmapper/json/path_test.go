package json

import (
	"encoding/json"
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/json/field"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolvePathValue(t *testing.T) {
	// Resolution of Old Trigger expression

	mapVal, _ := CoerceToObject("{\"myParam\":5}")
	path := ".myParam"
	newVal, err := ResolvePathValue(mapVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 5.0, newVal)

	// Resolution of Old Trigger expression
	arrVal, _ := CoerceToArray("[1,6,3]")
	path = "[1]"
	newVal, err = ResolvePathValue(arrVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 6.0, newVal)

	mapVal, _ = CoerceToObject("{\"myParam\":{\"nestedMap\":1}}")
	path = ".myParam.nestedMap"
	newVal, err = ResolvePathValue(mapVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 1.0, newVal)

	mapVal, _ = CoerceToObject("{\"myParam\":{\"nestedMap\":1}}")
	path = `["myParam"].nestedMap`
	newVal, err = ResolvePathValue(mapVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 1.0, newVal)

	mapVal, _ = CoerceToObject("{\"myParam\":{\"nestedMap\":1}}")
	path = `.myParam["nestedMap"]`
	newVal, err = ResolvePathValue(mapVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 1.0, newVal)

	arrVal, _ = CoerceToArray("[{\"nestedMap1\":1},{\"nestedMap2\":2}]")
	path = "[1].nestedMap2"
	newVal, err = ResolvePathValue(arrVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 2.0, newVal)

	mapVal, _ = CoerceToObject("{\"myParam\":{\"nestedArray\":[7,8,9]}}")
	path = ".myParam.nestedArray[1]"
	newVal, err = ResolvePathValue(mapVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 8.0, newVal)

	arrVal, _ = CoerceToArray("[{\"nestedMap1\":1},{\"nestedMap2\":{\"nestedArray\":[7,8,9]}}]")
	path = "[1].nestedMap2.nestedArray[2]"
	newVal, err = ResolvePathValue(arrVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 9.0, newVal)

	mapVal, _ = CoerceToObject("{\"myParam\":{\"nestedArray\":[7,8,9]}}")
	path = ".myParam.nestedArray"
	newVal, err = ResolvePathValue(mapVal, path)
	assert.Nil(t, err)
	//todo check if array

	arrVal, _ = CoerceToArray("[{\"nestedMap1\":1},{\"nestedMap2\":{\"nestedArray\":[7,8,9]}}]")
	path = "[1].nestedMap2"
	newVal, err = ResolvePathValue(arrVal, path)
	assert.Nil(t, err)
	//todo check if map

}

func TestPathSetValue(t *testing.T) {
	// Resolution of Old Trigger expression

	mapVal, _ := CoerceToObject("{\"myParam\":5}")
	path := ".myParam"
	v, err := PathSetValue(mapVal, path, 6)
	assert.Nil(t, err)
	newVal, err := ResolvePathValue(mapVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 6, newVal)

	// Resolution of Old Trigger expression
	arrVal, _ := CoerceToArray("[1,6,3]")
	path = "[1]"
	v, err = PathSetValue(arrVal, path, 4)
	assert.Nil(t, err)
	assert.Equal(t, 4, arrVal[1])
	newVal, err = ResolvePathValue(arrVal, path)
	assert.Nil(t, err)
	assert.Equal(t, float64(4), newVal)
	//
	mapVal, _ = CoerceToObject("{\"myParam\":{\"nestedMap\":1}}")
	path = ".myParam.nestedMap"
	assert.Nil(t, err)
	v, err = PathSetValue(mapVal, path, 7)
	newVal, err = ResolvePathValue(mapVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 7, newVal)

	mapVal, _ = CoerceToObject("{\"myParam\":{\"nestedMap\":1}}")
	path = `["myParam"].nestedMap`
	assert.Nil(t, err)
	v, err = PathSetValue(mapVal, path, 7)
	newVal, err = ResolvePathValue(mapVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 7, newVal)

	mapVal, _ = CoerceToObject("{\"myParam\":{\"nestedMap\":1}}")
	path = `.myParam["nestedMap"]`
	assert.Nil(t, err)
	v, err = PathSetValue(mapVal, path, 7)
	newVal, err = ResolvePathValue(mapVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 7, newVal)

	arrVal, _ = CoerceToArray("[{\"nestedMap1\":1},{\"nestedMap2\":2}]")
	path = "[1].nestedMap2"
	v, err = PathSetValue(arrVal, path, 3)
	assert.Nil(t, err)
	newVal, err = ResolvePathValue(v, path)
	assert.Nil(t, err)
	assert.Equal(t, float64(3), newVal)
	//
	mapVal, _ = CoerceToObject("{\"myParam\":{\"nestedArray\":[7,8,9]}}")
	path = ".myParam.nestedArray[1]"
	v, err = PathSetValue(mapVal, path, 1)
	assert.Nil(t, err)
	newVal, err = ResolvePathValue(mapVal, path)
	assert.Nil(t, err)
	assert.Equal(t, 1, newVal)

	arrVal, _ = CoerceToArray("[{\"nestedMap1\":1},{\"nestedMap2\":{\"nestedArray\":[7,8,9]}}]")
	path = "[1].nestedMap2.nestedArray[2]"
	v, err = PathSetValue(arrVal, path, 5)
	assert.Nil(t, err)
	fmt.Println(fmt.Sprintf("%+v", v))
	newVal, err = ResolvePathValue(v, path)
	assert.Nil(t, err)
	assert.Equal(t, float64(5), newVal)

	newVal, _ = CoerceToObject("{\"myParam\":{\"nestedArray\":[7,8,9]}}")
	path = ".myParam.nestedArray"
	v, err = PathSetValue(newVal, path, 3)
	assert.Nil(t, err)
	newVal, err = ResolvePathValue(v, path)
	assert.Nil(t, err)

	//todo check if array
	//arrVal,_ = CoerceToArray("[{\"nestedMap1\":1},{\"nestedMap2\":{\"nestedArray\":[7,8,9]}}]")
	//path = "[1].nestedMap2"
	//assert.Nil(t, err)
	//err = PathSetValue(arrVal, path, 3)
	//newVal,err = ResolvePathValue(arrVal, path)
	//assert.Nil(t, err)
	//////todo check if map
}

func PathSetValue(value interface{}, path string, attrValue interface{}) (interface{}, error) {
	mapField, err := field.ParseMappingField(path)
	if err != nil {
		return nil, err
	}

	return SetFieldValue(attrValue, value, mapField)
}

//Test purpose to duplicate this func to avoid circle import.(test only)
func CoerceToObject(v interface{}) (map[string]interface{}, error) {
	switch t := v.(type) {
	case string:
		m := make(map[string]interface{})
		if t != "" {
			err := json.Unmarshal([]byte(t), &m)
			if err != nil {
				return nil, fmt.Errorf("unable to coerce %#v to map[string]interface{}", t)
			}
		}
		return m, nil
	}
	return nil, nil
}

func CoerceToArray(v interface{}) ([]interface{}, error) {
	switch t := v.(type) {
	case string:
		a := make([]interface{}, 0)
		if t != "" {
			err := json.Unmarshal([]byte(t), &a)
			if err != nil {
				return nil, fmt.Errorf("unable to coerce %#v to map[string]interface{}", t)
			}
		}
		return a, nil
	}

	return nil, nil
}
