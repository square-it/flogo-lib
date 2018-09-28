package exprmapper

import (
	"testing"

	"fmt"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/stretchr/testify/assert"
)

func TestGetMapValueFunction(t *testing.T) {
	mappingValue := `string.concat("ddddd",$activity[a1].field.id)`
	v, err := GetExpresssionValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "dddddd", v)

	mappingValue = `string.concat("ddddd",$activity[a1].field.id, string.concat($activity[a1].field.id,$activity[a1].field.id))`
	v, err = GetExpresssionValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "dddddddd", v)

}

func TestGetMapValueExpression(t *testing.T) {
	mappingValue := `string.length(string.concat("ddddd",$activity[a1].field.id)) == 6`
	v, err := GetExpresssionValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, true, v)

	//
	mappingValue = `$activity[a1].field.id == "d"`
	v, err = GetExpresssionValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, true, v)

	mappingValue = `($activity[a1].field.id == "d") == true`
	v, err = GetExpresssionValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, true, v)

	mappingValue = `$activity[a1].field.id == "d" ? $activity[a1].field.id : string.concat("ssss",$activity[a1].field.id)`
	v, err = GetExpresssionValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "d", v)

	mappingValue = `$activity[a1].field.id == "d" ? string.concat("ssss",$activity[a1].field.id) : $activity[a1].field.id`
	v, err = GetExpresssionValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "ssssd", v)

	mappingValue = `$activity[a1].field.id != "d" ? string.concat("ssss",$activity[a1].field.id) : $activity[a1].field.id`
	v, err = GetExpresssionValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "d", v)

	mappingValue = `$activity[a1].field.id != "d" ? "dddd":"ssss"`
	v, err = GetExpresssionValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "ssss", v)

	mappingValue = `$activity[a1].field.id == "d" ? "dddd":"ssss"`
	v, err = GetExpresssionValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "dddd", v)

	mappingValue = ` 2>1 ? "dddd":"ssss"`
	v, err = GetExpresssionValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "dddd", v)

	mappingValue = ` 2<1 ? "dddd":"ssss"`
	v, err = GetExpresssionValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "ssss", v)

}

func TestMappingRef(t *testing.T) {
	mappingValue := `$activity[a1].field.id`
	v, err := GetExpresssionValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "d", v)

	mappingValue = `ddddd`
	v, err = GetExpresssionValue(mappingValue, GetSimpleScope("_A.a1.field", `{"id":"d"}`), GetTestResolver())
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "ddddd", v)

}

func GetSimpleScope(name, value string) data.Scope {
	a, _ := data.NewAttribute(name, data.TypeObject, value)
	maps := make(map[string]*data.Attribute)
	maps[name] = a
	scope := data.NewFixedScope(maps)
	scope.SetAttrValue(name, value)
	return scope
}

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
