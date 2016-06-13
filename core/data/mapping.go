package data

import (
	"fmt"
	"strconv"
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
	Type  MappingType `json:"type"`

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
// return error
func (m *Mapper) Apply(inputScope Scope, outputScope Scope) {

	//todo validate types
	for _, mapping := range m.mappings {

		switch mapping.Type {
		case MtAssign:

			var attrValue interface{}
			var exists bool

			attrName, attrPath, pathType := GetAttrPath(mapping.Value)

			attrValue, exists = inputScope.GetAttrValue(attrName)

			if exists && len(attrPath) > 0 {
				attrType, _ := inputScope.GetAttrType(attrName)
				if attrType == PARAMS {
					valMap := attrValue.(map[string]string)
					attrValue, exists = valMap[attrPath]
				} else if attrType == ARRAY && pathType == PT_ARRAY {
					//assigning part of array
					idx, _ := strconv.Atoi(attrPath)
					//todo handle err
					valArray := attrValue.([]interface{})
					attrValue = valArray[idx]
				} else {
					//for now assume if we have a path, attr is "object"
					valMap := attrValue.(map[string]interface{})
					attrValue, exists = valMap[attrPath]
				}
			}

			//todo implement type conversion
			if exists {

				attrName, attrPath, pathType := GetAttrPath(mapping.MapTo)
				toType, oe := outputScope.GetAttrType(attrName)

				if !oe {
					//todo handle attr dne
					fmt.Printf("Attr %s not found in output scope\n", attrName)
					return
				}

				switch pathType {
				case PT_SIMPLE:
					outputScope.SetAttrValue(mapping.MapTo, attrValue)
				case PT_ARRAY:
					val, _ := outputScope.GetAttrValue(attrName)
					if toType == ARRAY {
						var valArray []interface{}
						if val == nil {
							//what should we do in this case, construct the array?
							//valArray = make(map[string]string)
						} else {
							valArray = val.([]interface{})
						}

						idx, _ := strconv.Atoi(attrPath)
						//todo handle err
						valArray[idx] = attrValue

						outputScope.SetAttrValue(attrName, valArray)
					} else {
						//todo throw error.. not an ARRAY
					}
				case PT_MAP:
					val, _ := outputScope.GetAttrValue(attrName)
					if toType == PARAMS {
						var valMap map[string]string
						if val == nil {
							valMap = make(map[string]string)
						} else {
							valMap = val.(map[string]string)
						}
						strVal, _ := CoerceToString(attrValue)
						valMap[attrPath] = strVal

						outputScope.SetAttrValue(attrName, valMap)
					} else if toType == OBJECT {
						var valMap map[string]interface{}
						if val == nil {
							valMap = make(map[string]interface{})
						} else {
							valMap = val.(map[string]interface{})
						}
						valMap[attrPath] = attrValue

						outputScope.SetAttrValue(attrName, valMap)
					} else {
						//todo throw error.. not a OBJECT or PARAMS
					}
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
