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
		value, err := makeInterface(value)
		if err != nil {
			value = value
		}
		return value, nil
	}
	return GetFieldValue(value, mappingField)
}

func GetFieldValue(data interface{}, mappingField *field.MappingField) (interface{}, error) {
	var jsonParsed *Container
	var err error

	switch data.(type) {
	case string:
		jsonParsed, err = ParseJSON([]byte(data.(string)))
	case map[string]interface{}, map[string]string:
		jsonParsed, err = Consume(data)
	default:
		//Take is as string to handle
		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		jsonParsed, err = ParseJSON(b)
	}

	if err != nil {
		return nil, err

	}
	return handleGetValue(&JSONData{container: jsonParsed, rw: sync.RWMutex{}}, mappingField.Getfields())
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

func makeInterface(value interface{}) (interface{}, error) {

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
