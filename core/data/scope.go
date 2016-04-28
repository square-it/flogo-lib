package data

// Scope is an interface for getting and setting Attribute values
type Scope interface {

	// GetAttrType gets the type of the specified attribute
	GetAttrType(attrName string) (attrType string, exists bool)

	// GetAttrValue gets the value of the specified attribute
	GetAttrValue(attrName string) (value interface{}, exists bool)

	// SetAttrValue sets the value of the specified attribute
	SetAttrValue(attrName string, value interface{})
}

// SimpleScope is a basic implementation of a scope
type SimpleScope struct {
	parentScope Scope
	attrs       map[string]*Attribute
}

// NewSimpleScope creates a new SimpleScope
func NewSimpleScope(attrs []*Attribute, parentScope Scope) *SimpleScope {

	scope := &SimpleScope{
		parentScope: parentScope,
		attrs:       make(map[string]*Attribute),
	}

	for _, attr := range attrs {
		scope.attrs[attr.Name] = attr
	}

	return scope
}

// GetAttrType implements Scope.GetAttrType
func (s *SimpleScope) GetAttrType(attrName string) (attrType string, exists bool) {

	attr, found := s.attrs[attrName]

	if found {
		return attr.Type, true
	}

	if s.parentScope != nil {
		return s.parentScope.GetAttrType(attrName)
	}

	return "", false
}

// GetAttrValue implements Scope.GetAttrValue
func (s *SimpleScope) GetAttrValue(attrName string) (value interface{}, exists bool) {

	attr, found := s.attrs[attrName]

	if found {
		return attr.Value, true
	}

	if s.parentScope != nil {
		return s.parentScope.GetAttrValue(attrName)
	}

	return nil, false
}

// SetAttrValue implements Scope.SetAttrValue
func (s *SimpleScope) SetAttrValue(attrName string, value interface{}) {

	attr, found := s.attrs[attrName]

	if found {
		attr.Value = value
	}
	//todo return error? how do we determine type
}
