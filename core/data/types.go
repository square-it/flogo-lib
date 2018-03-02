package data

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Type denotes a data type
type Type int

const (
	TypeAny           Type = iota
	TypeString
	TypeInteger
	TypeNumber
	TypeBoolean
	TypeObject
	TypeArray
	TypeParams
	TypeComplexObject
)

var types = [...]string{
	"any",
	"string",
	"integer",
	"number",
	"boolean",
	"object",
	"array",
	"params",
	"any",
	"complex_object",
}

var typeMap = map[string]Type{
	"any":            TypeAny,
	"string":         TypeString,
	"integer":        TypeInteger,
	"number":         TypeNumber,
	"boolean":        TypeBoolean,
	"object":         TypeObject,
	"array":          TypeArray,
	"params":         TypeParams,
	"complex_object": TypeComplexObject,
}

func (t Type) String() string {
	return types[t]
}

// ToTypeEnum get the data type that corresponds to the specified name
func ToTypeEnum(typeStr string) (Type, bool) {

	dataType, found := typeMap[strings.ToLower(typeStr)]

	return dataType, found
}

// GetType get the Type of the supplied value
func GetType(val interface{}) (Type, error) {

	switch t := val.(type) {
	case string:
		return TypeString, nil
	case int:
		return TypeInteger, nil
	case float64:
		return TypeNumber, nil
	case json.Number:
		return TypeNumber, nil
	case bool:
		return TypeBoolean, nil
	case map[string]interface{}:
		return TypeObject, nil
	case []interface{}:
		return TypeArray, nil
	case ComplexObject:
		return TypeComplexObject, nil
	default:
		return TypeAny, fmt.Errorf("unable to determine type of %#v", t)
	}
}

func IsSimpleType(val interface{}) bool {

	switch val.(type) {
	case string, int, float64, json.Number, bool:
		return true
	default:
		return false
	}
}
