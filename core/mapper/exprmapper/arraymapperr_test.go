package exprmapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//7. array mapping with leaf field
func TestArrayMapping(t *testing.T) {
	mappingValue := `{
    "fields": [
        {
            "from": "$.street",
            "to": "$.street",
            "type": "primitive"
        },
        {
            "from": "$.zipcode",
            "to": "$.zipcode",
            "type": "primitive"
        },
        {
            "from": "$.state",
            "to": "$.state",
            "type": "primitive"
        }
    ],
    "from": "$activity[a1].field.addresses",
    "to": "field.addresses",
    "type": "foreach"
}`

	arrayData := `{
    "person": "name",
    "addresses": [
        {
            "street": "street",
            "zipcode": 77479,
            "state": "tx"
        }
    ]
}`

	array, err := ParseArrayMapping(mappingValue)
	assert.Nil(t, array.Validate())

	inputScope := GetObjectFieldScope("_A.a1.field", arrayData)
	outputScope := GetObjectFieldScope("field", "")
	err = array.DoArrayMapping(inputScope, outputScope, GetTestResolver())
	assert.Nil(t, err)

	arr, ok := outputScope.GetAttr("field")
	assert.True(t, ok)
	assert.Equal(t, "street", arr.Value().(map[string]interface{})["addresses"].([]interface{})[0].(map[string]interface{})["street"])
	assert.Equal(t, float64(77479), arr.Value().(map[string]interface{})["addresses"].([]interface{})[0].(map[string]interface{})["zipcode"])
	assert.Equal(t, "tx", arr.Value().(map[string]interface{})["addresses"].([]interface{})[0].(map[string]interface{})["state"])

}

//8. array mappping with static function and leaf field
func TestArrayMappingWithFunction(t *testing.T) {
	mappingValue := `{
    "fields": [
        {
            "from": "string.concat(\"this stree name: \", $.street)",
            "to": "$.street",
            "type": "primitive"
        },
        {
            "from": "string.concat(\"The zipcode is: \",$.zipcode)",
            "to": "$.zipcode",
            "type": "primitive"
        },
        {
            "from": "$.state",
            "to": "$.state",
            "type": "primitive"
        }
    ],
    "from": "$activity[a1].field.addresses",
    "to": "field.addresses",
    "type": "foreach"
}`

	arrayData := `{
    "person": "name",
    "addresses": [
        {
            "street": "street",
            "zipcode": 77479,
            "state": "tx"
        }
    ]
}`

	array, err := ParseArrayMapping(mappingValue)
	assert.Nil(t, err)
	assert.Nil(t, array.Validate())

	inputScope := GetObjectFieldScope("_A.a1.field", arrayData)
	outputScope := GetObjectFieldScope("field", "")
	err = array.DoArrayMapping(inputScope, outputScope, GetTestResolver())
	assert.Nil(t, err)

	arr, ok := outputScope.GetAttr("field")
	assert.True(t, ok)
	assert.Equal(t, "this stree name: street", arr.Value().(map[string]interface{})["addresses"].([]interface{})[0].(map[string]interface{})["street"])
	assert.Equal(t, "The zipcode is: 77479", arr.Value().(map[string]interface{})["addresses"].([]interface{})[0].(map[string]interface{})["zipcode"])
	assert.Equal(t, "tx", arr.Value().(map[string]interface{})["addresses"].([]interface{})[0].(map[string]interface{})["state"])

}

//9. array mapping with other activity output
func TestArrayMappingWithUpstreamingOutput(t *testing.T) {
	mappingValue := `{
    "fields": [
        {
            "from": "string.concat(\"this stree name: \", $activity[a1].field.person)",
            "to": "$.street",
            "type": "primitive"
        },
        {
            "from": "string.concat(\"The zipcode is: \",$.zipcode)",
            "to": "$.zipcode",
            "type": "primitive"
        },
        {
            "from": "$.state",
            "to": "$.state",
            "type": "primitive"
        }
    ],
    "from": "$activity[a1].field.addresses",
    "to": "field.addresses",
    "type": "foreach"
}`

	arrayData := `{
    "person": "name",
    "addresses": [
        {
            "street": "street",
            "zipcode": 77479,
            "state": "tx"
        }
    ]
}`

	array, err := ParseArrayMapping(mappingValue)
	assert.Nil(t, err)
	assert.Nil(t, array.Validate())

	inputScope := GetObjectFieldScope("_A.a1.field", arrayData)
	outputScope := GetObjectFieldScope("field", "")
	err = array.DoArrayMapping(inputScope, outputScope, GetTestResolver())
	assert.Nil(t, err)

	arr, ok := outputScope.GetAttr("field")
	assert.True(t, ok)
	assert.Equal(t, "this stree name: name", arr.Value().(map[string]interface{})["addresses"].([]interface{})[0].(map[string]interface{})["street"])
	assert.Equal(t, "The zipcode is: 77479", arr.Value().(map[string]interface{})["addresses"].([]interface{})[0].(map[string]interface{})["zipcode"])
	assert.Equal(t, "tx", arr.Value().(map[string]interface{})["addresses"].([]interface{})[0].(map[string]interface{})["state"])

}
