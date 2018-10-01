package json

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/json/field"
	"strconv"
	"strings"

	"sync"

	"encoding/json"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("json")

func ResolvePathValue(value interface{}, refPath string) (interface{}, error) {
	mappingField, err := field.ParseMappingField(refPath)
	if err != nil {
		return nil, fmt.Errorf("parse mapping path [%s] failed, due to %s", err.Error())
	}

	if mappingField == nil || len(mappingField.Getfields()) <= 0 {
		value, err := toInterface(value)
		if err != nil {
			value = value
		}
		return value, nil
	}
	return GetFieldValue(value, mappingField)
}

func toInterface(value interface{}) (interface{}, error) {

	var paramMap interface{}

	if value == nil {
		return paramMap, nil
	}

	switch t := value.(type) {
	case string:
		err := json.Unmarshal([]byte(t), &paramMap)
		if err != nil {
			return nil, err
		}
		return paramMap, nil
	default:
		return value, nil
	}
	return paramMap, nil
}

func GetFieldValue(data interface{}, mappingField *field.MappingField) (interface{}, error) {
	var jsonParsed *Container
	var err error
	value, ok := data.(string)
	if ok {
		jsonParsed, err = ParseJSON([]byte(value))
	} else {
		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		jsonParsed, err = ParseJSON(b)
	}

	if err != nil {
		return nil, err

	}
	return getFieldValue(&JSONData{container: jsonParsed, rw: sync.RWMutex{}}, mappingField)
}

func getFieldValue(jsonData *JSONData, mappingField *field.MappingField) (interface{}, error) {
	return handleGetSpecialFields(jsonData, mappingField.Getfields())
}

func getRestArrayFieldName(fieldName string) string {
	if strings.Index(fieldName, "]") >= 0 {
		closeBracketIndex := strings.Index(fieldName, "]")
		if len(fieldName) == closeBracketIndex+1 {
			return ""
		}
		return fieldName[closeBracketIndex+2:]
	}
	return fieldName
}

func getFieldName(fieldName string) string {
	if strings.Index(fieldName, "[") >= 0 {
		return fieldName[0:strings.Index(fieldName, "[")]
	}

	return fieldName
}

func getFieldSliceIndex(fieldName string) (int, error) {
	if strings.Index(fieldName, "[") >= 0 {
		index := fieldName[strings.Index(fieldName, "[")+1 : strings.Index(fieldName, "]")]
		i, err := strconv.Atoi(index)

		if err != nil {
			return -2, nil
		}
		return i, nil
	}

	return -1, nil
}

func getNameInsideBrancket(fieldName string) string {
	if strings.Index(fieldName, "[") >= 0 {
		index := fieldName[strings.Index(fieldName, "[")+1 : strings.Index(fieldName, "]")]
		return index
	}

	return ""
}
