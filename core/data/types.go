package data

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Type denotes a data type
type Type int

const (
	ANY            Type = iota
	STRING
	INTEGER
	NUMBER
	BOOLEAN
	OBJECT
	ARRAY
	PARAMS
	COMPLEX_OBJECT
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
	"any":            ANY,
	"string":         STRING,
	"integer":        INTEGER,
	"number":         NUMBER,
	"boolean":        BOOLEAN,
	"object":         OBJECT,
	"array":          ARRAY,
	"params":         PARAMS,
	"complex_object": COMPLEX_OBJECT,
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
		return STRING, nil
	case int:
		return INTEGER, nil
	case float64:
		return NUMBER, nil
	case json.Number:
		return NUMBER, nil
	case bool:
		return BOOLEAN, nil
	case map[string]interface{}:
		return OBJECT, nil
	case []interface{}:
		return ARRAY, nil
	case ComplexObject:
		return COMPLEX_OBJECT, nil
	default:
		return ANY, fmt.Errorf("unable to determine type of %#v", t)
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
