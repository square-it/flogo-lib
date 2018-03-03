package util

import (
	"encoding/json"
	"fmt"
	"strconv"
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

func ConvertToInt(val interface{}) (int, error) {
	switch t := val.(type) {
	case int:
		return t, nil
	case int64:
		return int(t), nil
	case float64:
		return int(t), nil
	case json.Number:
		i, err := t.Int64()
		return int(i), err
	case string:
		return strconv.Atoi(t)
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

func ConvertToInterface(value interface{}) (interface{}, error) {

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

func ConvertToBytes(value interface{}) ([]byte, error) {
	if value != nil {
		switch t := value.(type) {
		case []byte:
			return t, nil
		case string:
			return []byte(t), nil
		default:
			return json.Marshal(value)
		}
	}
	return nil, fmt.Errorf("Cannot convert nil to []byte")
}

func ConvertToBool(value interface{}) (bool, error) {
	if value != nil {
		switch t := value.(type) {
		case bool:
			return t, nil
		case string:
			return strconv.ParseBool(t)
		default:
			str, err := ConvertToString(value)
			if err != nil {
				return false, err
			}
			return strconv.ParseBool(str)
		}
	}
	return false, nil
}
