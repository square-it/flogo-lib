package data

// Attribute is a simple structure used to define a data Attribute/property
type Attribute struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value interface{} `json:"value"`
}

type ProtoAttr struct {
	Attribute

	valueSet bool
}
