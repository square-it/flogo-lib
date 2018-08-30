package exprmapper

import (
	"testing"

	"fmt"
	"strings"

	"encoding/json"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/stretchr/testify/assert"
)

//1. activity mapping
func TestActivityMapping(t *testing.T) {
	mappingValue := `$activity[a1].field.id`
	v, err := GetMappingValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "d", v)
}

func TestMappingRef(t *testing.T) {
	mappingValue := `$activity[a1].field.id`
	v, err := GetMappingValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "d", v)

	mappingValue = `ddddd`
	v, err = GetMappingValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "ddddd", v)

}

//2. flow mapping
func TestFlowMapping(t *testing.T) {
	mappingValue := `$.field.id`
	v, err := GetMappingValue(mappingValue, GetSimpleScope("field", `{"id":"d"}`), data.GetBasicResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "d", v)
}

//3. function
func TestGetMapValueFunction(t *testing.T) {
	mappingValue := `string.concat("ddddd",$activity[a1].field.id)`
	v, err := GetMappingValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "dddddd", v)

	mappingValue = `string.concat("ddddd",$activity[a1].field.id, string.concat($activity[a1].field.id,$activity[a1].field.id))`
	v, err = GetMappingValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "dddddddd", v)
}

//4. function with different type
func TestGetMapValueFunctionTypes(t *testing.T) {
	mappingValue := `string.concat(123,$activity[a1].field.id)`
	v, err := GetMappingValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "123d", v)

	mappingValue = `string.concat("s-",$activity[a1].field.id, string.concat($activity[a1].field.id,$activity[a1].field.id),true)`
	v, err = GetMappingValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "s-dddtrue", v)

	mappingValue = `string.concat(123,$activity[a1].field.id, "dddd", 450)`
	v, err = GetMappingValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "123ddddd450", v)
}

//5. expression
//6. ternary expression
func TestGetMapValueExpression(t *testing.T) {
	mappingValue := `string.length(string.concat("ddddd",$activity[a1].field.id)) == 6`
	v, err := GetMappingValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, true, v)

	//
	mappingValue = `$activity[a1].field.id == "d"`
	v, err = GetMappingValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, true, v)

	mappingValue = `($activity[a1].field.id == "d") == true`
	v, err = GetMappingValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, true, v)

	mappingValue = `$activity[a1].field.id == "d" ? $activity[a1].field.id : string.concat("ssss",$activity[a1].field.id)`
	v, err = GetMappingValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "d", v)

	mappingValue = `$activity[a1].field.id == "d" ? string.concat("ssss",$activity[a1].field.id) : $activity[a1].field.id`
	v, err = GetMappingValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "ssssd", v)

	mappingValue = `$activity[a1].field.id != "d" ? string.concat("ssss",$activity[a1].field.id) : $activity[a1].field.id`
	v, err = GetMappingValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "d", v)

	mappingValue = `$activity[a1].field.id != "d" ? "dddd":"ssss"`
	v, err = GetMappingValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "ssss", v)

	mappingValue = `$activity[a1].field.id == "d" ? "dddd":"ssss"`
	v, err = GetMappingValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "dddd", v)

	mappingValue = ` 2>1 ? "dddd":"ssss"`
	v, err = GetMappingValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "dddd", v)

	mappingValue = ` 2<1 ? "dddd":"ssss"`
	v, err = GetMappingValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "ssss", v)

}

func GetSimpleScope(name, value string) data.Scope {
	a, _ := data.NewAttribute(name, data.TypeObject, value)
	maps := make(map[string]*data.Attribute)
	maps[name] = a
	scope := data.NewFixedScope(maps)
	scope.SetAttrValue(name, value)
	return scope
}

//7. array mapping
func TestArrayMapping(t *testing.T) {
	mappingValue := `{
    "fields": [
        {
            "from": "$.street",
            "to": "$$['street']",
            "type": "primitive"
        },
        {
            "from": "$.zipcode",
            "to": "$$['zipcode']",
            "type": "primitive"
        },
        {
            "from": "$.state",
            "to": "$$['state']",
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

	inputScope := GetSimpleScope("_A.a1.field", arrayData)
	outputScope := GetSimpleScope("field", "")
	err = array.DoArrayMapping(inputScope, outputScope, GetTestResolver())
	assert.Nil(t, err)

	arr, ok := outputScope.GetAttr("field")
	assert.True(t, ok)
	v, _ := json.Marshal(arr.Value().(map[string]interface{})["addresses"].([]interface{})[0])
	fmt.Println(string(v))
	assert.Equal(t, "street", arr.Value().(map[string]interface{})["addresses"].([]interface{})[0].(map[string]interface{})["street"])
	assert.Equal(t, float64(77479), arr.Value().(map[string]interface{})["addresses"].([]interface{})[0].(map[string]interface{})["zipcode"])
	assert.Equal(t, "tx", arr.Value().(map[string]interface{})["addresses"].([]interface{})[0].(map[string]interface{})["state"])

}

//8. array mapping with leaf field
//9. array mappping with static function and leaf field
//10. array mapping with other activity output

//For test  purpuse copy it from flogo-contrib
var resolver = &TestResolver{}

func GetTestResolver() data.Resolver {
	return resolver
}

type TestResolver struct {
}

func (r *TestResolver) Resolve(toResolve string, scope data.Scope) (value interface{}, err error) {

	var details *data.ResolutionDetails

	if strings.HasPrefix(toResolve, "${") {
		details, err = data.GetResolutionDetailsOld(toResolve)
	} else if strings.HasPrefix(toResolve, "$") {
		details, err = data.GetResolutionDetails(toResolve[1:])
	} else {
		return data.SimpleScopeResolve(toResolve, scope)
	}

	if err != nil {
		return nil, err
	}

	if details == nil {
		return nil, fmt.Errorf("unable to determine resolver for %s", toResolve)
	}

	switch details.ResolverName {
	case "activity":
		attr, exists := scope.GetAttr("_A." + details.Item + "." + details.Property)
		if !exists {
			return nil, fmt.Errorf("failed to resolve activity attr: '%s', not found in flow", details.Property)
		}
		value = attr.Value()
	default:
		return nil, fmt.Errorf("unsupported resolver: %s", details.ResolverName)
	}

	if details.Path != "" {
		value, err = data.PathGetValue(value, details.Path)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
	}

	return value, nil
}
