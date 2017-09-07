package mapper

import (
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"fmt"
)

type Factory interface {
	// NewMapper creates a new data.Mapper from the specified data.MapperDef
	NewMapper(mapperDef *data.MapperDef) data.Mapper

	// NewUniqueMapper creates a unique data.Mapper from the specified data.MapperDef
	// the ID can be used to facilitate use precompiled mappers
	NewUniqueMapper(ID string, mapperDef *data.MapperDef) data.Mapper
}

var factory Factory

func SetFactory(factory Factory) {
	factory = factory
}

func GetFactory() Factory {

	if factory == nil {
		factory = &BasicMapperFactory{}
	}

	return factory
}

type BasicMapperFactory struct {
}

func (mf *BasicMapperFactory) NewMapper(mapperDef *data.MapperDef) data.Mapper {
	return NewBasicMapper(mapperDef)
}

func (mf *BasicMapperFactory) NewUniqueMapper(ID string, mapperDef *data.MapperDef) data.Mapper {
	return NewBasicMapper(mapperDef)
}

// BasicMapper is a simple object holding and executing mappings
type BasicMapper struct {
	mappings []*data.MappingDef
}

// NewBasicMapper creates a new BasicMapper with the specified mappings
func NewBasicMapper(mapperDef *data.MapperDef) data.Mapper {

	var mapper BasicMapper
	mapper.mappings = mapperDef.Mappings

	return &mapper
}

// Mappings gets the mappings of the BasicMapper
func (m *BasicMapper) Mappings() []*data.MappingDef {
	return m.mappings
}

// Apply executes the mappings using the values from the input scope
// and puts the results in the output scope
//
// return error
func (m *BasicMapper) Apply(inputScope data.Scope, outputScope data.Scope) error {

	//todo validate types
	for _, mapping := range m.mappings {

		switch mapping.Type {
		case data.MtAssign:

			exprRep, ok := mapping.Value.(string)
			if !ok {
				return fmt.Errorf("Invalid assign: %v", mapping.Value)
			}

			lookupExpr := NewLookupExpr(exprRep)

			val, err := lookupExpr.Eval(inputScope)

			if err != nil {
				return err
			}

			assignExpr := NewAssignExpr(mapping.MapTo, val)

			_, err = assignExpr.Eval(outputScope)

			if err != nil {
				return err
			}

		case data.MtLiteral:
			assignExpr := NewAssignExpr(mapping.MapTo, mapping.Value)

			_, err := assignExpr.Eval(outputScope)

			if err != nil {
				return err
			}
		case data.MtExpression:
			//todo implement script mapping
		}
	}

	return nil
}
