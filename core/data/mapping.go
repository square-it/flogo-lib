package data

// MappingType is an enum for possible Mapping Types
type MappingType int

const (
	// AssignMT denotes an attribute to attribute assignment
	MtAssign MappingType = 1

	// LiteralMT denotes a literal to attribute assignment
	MtLiteral MappingType = 2

	// ExpressionMT denotes a expression execution to perform mapping
	MtExpression MappingType = 3
)

// Mapping is a simple structure that defines a mapping
type Mapping struct {
	//Type the mapping type
	Type MappingType `json:"type"`

	//Value the mapping value to execute to determine the result (lhs)
	Value string `json:"value"`

	//Result the name of attribute to place the result of the mapping in (rhs)
	MapTo string `json:"mapTo"`
}

// Mapper is a simple object holding and executing mappings
type Mapper struct {
	mappings []*Mapping
}

// NewMapper creates a new Mapper with the specified mappings
func NewMapper(mappings []*Mapping) *Mapper {

	var mapper Mapper
	mapper.mappings = mappings

	return &mapper
}

// Mappings gets the mappings of the mapper
func (m *Mapper) Mappings() []*Mapping {
	return m.mappings
}

// Apply executes the mappings using the values from the input scope
// and puts the results in the output scope
//
// todo: does the engine have to facilitate type conversion?
func (m *Mapper) Apply(inputScope Scope, outputScope Scope) {

	//todo validate types
	for _, mapping := range m.mappings {

		switch mapping.Type {
		case MtAssign:
			attrValue, exists := inputScope.GetAttrValue(mapping.Value)
			if exists {
				outputScope.SetAttrValue(mapping.MapTo, attrValue)
			}
		//todo: should we ignore if DNE - if we have to add dynamically what type do we use
		case MtLiteral:
			outputScope.SetAttrValue(mapping.MapTo, mapping.Value)
		case MtExpression:
			//todo implement script mapping
		}
	}
}
