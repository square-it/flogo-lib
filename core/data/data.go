package data

// Attribute is a simple structure used to define a data Attribute/property
type Attribute struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"` //todo: should be interface{}
}

// Scope is an interface for getting and setting Attribute values
type Scope interface {

	// GetAttrType gets the type of the specified attribute
	GetAttrType(attrName string) (attrType string, exists bool)

	// GetAttrValue gets the value of the specified attribute
	GetAttrValue(attrName string) (value string, exists bool)

	// SetAttrValue sets the value of the specified attribute
	SetAttrValue(attrName string, value string)
}
