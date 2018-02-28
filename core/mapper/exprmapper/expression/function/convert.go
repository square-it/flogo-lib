package function

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/util"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/witype"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/ref"
)

func ConvertToValue(value interface{}, dataType witype.TYPE) (interface{}, error) {

	var coerced interface{}
	var err error

	switch dataType {
	case witype.STRING:
		coerced, err = util.ConvertToString(value)
	case witype.INT64:
		coerced, err = util.ConvertToInt64(value)
	case witype.FLOAT:
		coerced, err = util.ConvertToFloat(value)
	case witype.FUNCTION:
		coerced, err = ConvertToFunction(value)
	case witype.REF:
		coerced, err = ConvertToRef(value)
	case witype.EXPRESSION:
		return value, nil
	}

	if err != nil {
		return nil, err
	}

	return coerced, nil
}

func ConvertToRef(val interface{}) (*ref.MappingRef, error) {

	logrus.Infof("Convert to ref type %s value %+v", reflect.TypeOf(val), val)
	switch val.(type) {
	case string:
		return ref.NewMappingRef(val.(string)), nil
	case *string:
		return ref.NewMappingRef(*val.(*string)), nil
	case *ref.MappingRef:
		return val.(*ref.MappingRef), nil
	case interface{}:
		v, err := json.Marshal(val)
		mapRef := &ref.MappingRef{}
		err = json.Unmarshal(v, mapRef)
		return mapRef, err
	}

	return nil, errors.New("Cannot convert to mapping ref")
}

func ConvertToFunction(val interface{}) (*FunctionExp, error) {
	if val == nil {
		return nil, nil
	}
	switch t := val.(type) {
	case *FunctionExp:
		return t, nil
	case string:
		logrus.Debug("Convert function from string.")
		function := &FunctionExp{}
		err := json.Unmarshal(val.([]byte), function)
		if err != nil {
			return nil, err
		}
		return function, nil
	case map[string]interface{}:
		v, err := json.Marshal(val)
		if err != nil {
			return nil, err
		}
		function := &FunctionExp{}
		err = json.Unmarshal(v, function)
		if err != nil {
			return nil, err
		}
		return function, nil
	case interface{}:
		v, err := json.Marshal(val)
		if err != nil {
			return nil, err
		}
		function := &FunctionExp{}
		err = json.Unmarshal(v, function)
		if err != nil {
			return nil, err
		}
		return function, nil
	default:
		logrus.Debugf("Convert function from type %s", reflect.TypeOf(val))
		return nil, fmt.Errorf("Unable to Convert %#v to function", t)
	}
}
