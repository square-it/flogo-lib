package data

import "encoding/json"

// Attribute is a simple structure used to define a data Attribute/property
type Attribute struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// TypedValue is a value with a type
type TypedValue struct {
	Type  Type        `json:"type"`
	Value interface{} `json:"value"`
}

// MarshalJSON implements json.Marshaler.MarshalJSON
func (tv *TypedValue) MarshalJSON() ([]byte, error) {

	return json.Marshal(&struct {
		Type  string      `json:"type"`
		Value interface{} `json:"value"`
	}{
		Type:  tv.Type.String(),
		Value: tv.Value,
	})
}

// UnmarshalJSON implements json.Unmarshaler.UnmarshalJSON
func (tv *TypedValue) UnmarshalJSON(data []byte) error {

	ser := &struct {
		Type  string      `json:"type"`
		Value interface{} `json:"value"`
	}{}

	if err := json.Unmarshal(data, ser); err != nil {
		return err
	}

	tv.Type, _ = ToTypeEnum(ser.Type)
	tv.Value = ser.Value

	return nil
}
