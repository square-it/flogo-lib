package util

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func ConvertToString(val interface{}) (string, error) {
	switch t := val.(type) {
	case string:
		return t, nil
	case int64:
		return strconv.FormatInt(t, 10), nil
	case int:
		return strconv.Itoa(t), nil
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64), nil
	case json.Number:
		return t.String(), nil
	case bool:
		return strconv.FormatBool(t), nil
	case nil:
		return "", nil
	case map[string]interface{}:
		b, err := json.Marshal(t)
		if err != nil {
			return "", err
		}
		return string(b), nil
	default:
		return "", fmt.Errorf("Unable to Coerce %#v to string", t)
	}
}

func ConvertToInt64(val interface{}) (int64, error) {
	switch t := val.(type) {
	case int:
		return int64(t), nil
	case int64:
		return t, nil
	case float64:
		return int64(t), nil
	case json.Number:
		i, err := t.Int64()
		return int64(i), err
	case string:
		return strconv.ParseInt(t, 10, 64)
	case bool:
		if t {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("Unable to coerce %#v to int64", val)
	}
}

func ConvertToFloat(val interface{}) (float64, error) {
	switch t := val.(type) {
	case int:
		return float64(t), nil
	case int64:
		return float64(t), nil
	case float64:
		return t, nil
	case json.Number:
		return t.Float64()
	case string:
		return strconv.ParseFloat(t, 64)
	case bool:
		if t {
			return 1.0, nil
		}
		return 0.0, nil
	case nil:
		return 0.0, nil
	default:
		return 0.0, fmt.Errorf("Unable to coerce %#v to float64", val)
	}
}

// ConvertToBoolean coerce a value to a boolean
func ConvertToBoolean(val interface{}) (bool, error) {
	switch t := val.(type) {
	case bool:
		return t, nil
	case int:
		return t != 0, nil
	case int64:
		return t != 0, nil
	case float64:
		return t != 0.0, nil
	case json.Number:
		i, err := t.Int64()
		return i != 0, err
	case string:
		return strconv.ParseBool(t)
	case nil:
		return false, nil
	default:
		return false, fmt.Errorf("Unable to coerce %#v to bool", val)
	}
}

func ConvertToObject(val interface{}) (map[string]interface{}, error) {

	switch t := val.(type) {
	case map[string]interface{}:
		return t, nil
	default:
		return nil, fmt.Errorf("Unable to coerce %#v to map[string]interface{}", val)
	}
}

func ConvertToArray(val interface{}) ([]interface{}, error) {

	switch t := val.(type) {
	case []interface{}:
		return t, nil
	case []map[string]interface{}:
		var a []interface{}
		for _, v := range t {
			a = append(a, v)
		}
		return a, nil
	default:
		return nil, fmt.Errorf("Unable to coerce %#v to []interface{}", val)
	}
}

func ConvertToAny(val interface{}) (interface{}, error) {

	switch t := val.(type) {

	case json.Number:

		if strings.Contains(t.String(), ".") {
			return t.Float64()
		} else {
			return t.Int64()
		}
	default:
		return val, nil
	}
}

func ConvertToParams(val interface{}) (map[string]string, error) {

	switch t := val.(type) {
	case map[string]string:
		return t, nil
	case map[string]interface{}:

		var m = make(map[string]string, len(t))
		for k, v := range t {

			mVal, err := ConvertToString(v)
			if err != nil {
				return nil, err
			}
			m[k] = mVal
		}
		return m, nil
	case map[interface{}]string:

		var m = make(map[string]string, len(t))
		for k, v := range t {

			mKey, err := ConvertToString(k)
			if err != nil {
				return nil, err
			}
			m[mKey] = v
		}
		return m, nil
	case map[interface{}]interface{}:

		var m = make(map[string]string, len(t))
		for k, v := range t {

			mKey, err := ConvertToString(k)
			if err != nil {
				return nil, err
			}

			mVal, err := ConvertToString(v)
			if err != nil {
				return nil, err
			}
			m[mKey] = mVal
		}
		return m, nil
	default:
		return nil, fmt.Errorf("Unable to coerce %#v to map[string]string", val)
	}
}
