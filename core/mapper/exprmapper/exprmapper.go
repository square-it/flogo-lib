package exprmapper

import (
	"errors"
	"fmt"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/json/field"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/json"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/ref"
	"github.com/TIBCOSoftware/flogo-lib/logger"

	//Pre registry all function for now
	_ "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/array/length"
	_ "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/number/random"
	_ "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/concat"
	_ "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/equals"
	_ "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/equalsignorecase"
	_ "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/length"
	_ "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/function/string/substring"
)

var log = logger.GetLogger("mapper")

const (
	MAP_TO_INPUT = "$INPUT"
)

func MapExpreesion(mapping *data.MappingDef, inputScope, outputScope data.Scope, resolver data.Resolver) error {
	mappingValue, err := GetExpresssionValue(mapping.Value, inputScope, resolver)
	if err != nil {
		return err
	}
	err = SetValueToOutputScope(mapping.MapTo, outputScope, mappingValue)
	if err != nil {
		err = fmt.Errorf("Set value %+v to output [%s] error - %s", mappingValue, mapping.MapTo, err.Error())
		log.Error(err)
		return err
	}
	log.Debugf("Set value %+v to %s Done", mappingValue, mapping.MapTo)
	return nil
}

func MapAssign(mapping *data.MappingDef, inputScope, outputScope data.Scope, resolver data.Resolver) error {
	mappingValue, err := GetMappingValue(mapping.Value, inputScope, resolver)
	if err != nil {
		return err
	}
	err = SetValueToOutputScope(mapping.MapTo, outputScope, mappingValue)
	if err != nil {
		err = fmt.Errorf("Set value %+v to output [%s] error - %s", mappingValue, mapping.MapTo, err.Error())
		log.Error(err)
		return err
	}
	log.Debugf("Set value %+v to %s Done", mappingValue, mapping.MapTo)
	return nil
}

func GetExpresssionValue(mappingV interface{}, inputScope data.Scope, resolver data.Resolver) (interface{}, error) {
	mappingValue, ok := mappingV.(string)
	if !ok {
		return mappingV, nil
	}
	exp, err := expression.ParseExpression(mappingValue)
	if err == nil {
		//flogo expression
		expValue, err := exp.EvalWithScope(inputScope, resolver)
		if err != nil {
			return nil, fmt.Errorf("Execution failed for mapping [%s] due to error - %s", mappingValue, err.Error())
		}
		return expValue, nil
	} else {
		return GetMappingValue(mappingV, inputScope, resolver)
	}
}

func GetMappingValue(mappingV interface{}, inputScope data.Scope, resolver data.Resolver) (interface{}, error) {
	mappingValue, ok := mappingV.(string)
	if !ok {
		return mappingV, nil
	}
	if !isMappingRef(mappingValue) {
		log.Debugf("Mapping value is literal set directly to field")
		log.Debugf("Mapping ref %s and value %+v", mappingValue, mappingValue)
		return mappingValue, nil
	} else {
		mappingref := ref.NewMappingRef(mappingValue)
		mappingValue, err := mappingref.GetValue(inputScope, resolver)
		if err != nil {
			return nil, fmt.Errorf("Get value from ref [%s] error - %s", mappingref.GetRef(), err.Error())

		}
		log.Debugf("Mapping ref %s and value %+v", mappingValue, mappingValue)
		return mappingValue, nil
	}
	return nil, nil
}

func SetValueToOutputScope(mapTo string, outputScope data.Scope, value interface{}) error {
	mapField, err := field.ParseMappingField(mapTo)
	if err != nil {
		return err
	}

	actRootField, err := ref.GetMapToAttrName(mapField)
	if err != nil {
		return err
	}

	fields := mapField.Getfields()
	if len(fields) == 1 && !ref.HasArray(fields[0]) {
		//No complex mapping exist
		return SetAttribute(actRootField, value, outputScope)
	} else if ref.HasArray(fields[0]) || len(fields) > 1 {
		//Complex mapping
		return settValueToComplexObject(mapField, actRootField, outputScope, value)
	} else {
		return fmt.Errorf("No field name found for mapTo [%s]", mapTo)
	}

}

func settValueToComplexObject(mapField *field.MappingField, fieldName string, outputScope data.Scope, value interface{}) error {
	complexVlaueIn, err := ref.GetValueFromOutputScope(mapField, outputScope)
	if err != nil {
		return err
	}
	pathfields, err := ref.GetMapToPathFields(mapField)
	if err != nil {
		return err
	}

	log.Debugf("Set value %+v to fields %s", value, pathfields)
	complexValue, err2 := json.SetFieldValue(value, complexVlaueIn, pathfields)
	if err2 != nil {
		return err2
	}

	return SetAttribute(fieldName, complexValue, outputScope)
}

func isMappingRef(mappingref string) bool {
	if mappingref == "" || !strings.HasPrefix(mappingref, "$") {
		return false
	}
	return true
}

func SetAttribute(fieldName string, value interface{}, outputScope data.Scope) error {
	//Set Attribute value back to attribute
	attribute, exist := outputScope.GetAttr(fieldName)
	if exist {
		switch attribute.Type() {
		case data.TypeComplexObject:
			complexObject := attribute.Value().(*data.ComplexObject)
			newComplexObject := &data.ComplexObject{Metadata: complexObject.Metadata, Value: value}
			outputScope.SetAttrValue(fieldName, newComplexObject)
		default:
			outputScope.SetAttrValue(fieldName, value)
		}

	} else {
		return errors.New("Cannot found attribute " + fieldName + " at output scope")
	}
	return nil
}

func RemovePrefixInput(str string) string {
	if str != "" && strings.HasPrefix(str, MAP_TO_INPUT) {
		//Remove $INPUT for mapTo
		newMapTo := str[len(MAP_TO_INPUT):]
		if strings.HasPrefix(newMapTo, ".") {
			newMapTo = newMapTo[1:]
		}
		str = newMapTo
	}
	return str
}
