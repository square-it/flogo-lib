package data

import (
	"fmt"
	"strings"
)

// MappingType is an enum for possible Mapping Types
type MappingType int

const (
	// MtAssign denotes an attribute to attribute assignment
	MtAssign MappingType = 1

	// MtLiteral denotes a literal to attribute assignment
	MtLiteral MappingType = 2

	// MtExpression denotes a expression execution to perform mapping
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

			var attrValue interface{}
			var exists bool

			inAttrName := mapping.Value

			//todo move this code
			if inAttrName[0] == '[' {

				if inAttrName[len(inAttrName)-1] != ']' {
					idx := strings.Index(inAttrName, "]")

					mapAttrName := inAttrName[idx+2:]

					fmt.Printf("AttrName: %s\n", inAttrName[:idx+1])
					fmt.Printf("MapAttrName: %s\n", mapAttrName)

					val, attrExists := inputScope.GetAttrValue(inAttrName[:idx+1])
					fmt.Printf("val: %v\n", val)

					if attrExists {

						valMap := val.(map[string]interface{})
						attrValue, exists = valMap[mapAttrName]
					}
				} else {
					attrValue, exists = inputScope.GetAttrValue(mapping.Value)
				}
			} else {
				attrValue, exists = inputScope.GetAttrValue(mapping.Value)
			}

			//idx := strings.Index(mapping.MapTo,".")

			//todo implement type conversion
			if exists {

				idx := strings.Index(mapping.MapTo, ".")

				if idx > -1 {
					attrName := mapping.MapTo[:idx]
					mapAttrName := mapping.MapTo[idx+1:]

					//assigning to map value
					toType, oe := outputScope.GetAttrType(attrName)

					if oe {
						if toType == PARAMS.String() {
							val, _ := outputScope.GetAttrValue(attrName)
							var valMap map[string]string
							if val == nil {
								valMap = make(map[string]string)
							} else {
								valMap = val.(map[string]string)
							}
							strVal, _ := CoerceToString(attrValue)
							valMap[mapAttrName] = strVal

							outputScope.SetAttrValue(attrName, valMap)

						} else {

							//error, not a map (or object?)
						}
					} else {
						fmt.Printf("Attr %s not found in output scope\n", attrName)
					}

				} else {
					outputScope.SetAttrValue(mapping.MapTo, attrValue)
				}
			}
		//todo: should we ignore if DNE - if we have to add dynamically what type do we use
		case MtLiteral:
			outputScope.SetAttrValue(mapping.MapTo, mapping.Value)
		case MtExpression:
			//todo implement script mapping
		}
	}
}
