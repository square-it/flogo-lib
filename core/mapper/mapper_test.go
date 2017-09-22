package mapper

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

func TestLiteralMapper(t *testing.T) {

	factory := GetFactory()

	mapping1 := &data.MappingDef{Type: data.MtLiteral, Value: "1", MapTo: "Simple"}
	mapping2 := &data.MappingDef{Type: data.MtLiteral, Value: 2, MapTo: "Obj.key"}
	mapping3 := &data.MappingDef{Type: data.MtLiteral, Value: 3, MapTo: "Array[2]"}
	mapping4 := &data.MappingDef{Type: data.MtLiteral, Value: "4", MapTo: "Params.paramKey"}

	mappings := []*data.MappingDef{mapping1, mapping2, mapping3, mapping4}

	mapper := factory.NewMapper(&data.MapperDef{Mappings: mappings})

	attr1 := &data.Attribute{Name: "Simple", Type: data.INTEGER}
	attr2 := &data.Attribute{Name: "Obj", Type: data.OBJECT}
	attr3 := &data.Attribute{Name: "Array", Type: data.ARRAY}
	attr4 := &data.Attribute{Name: "Params", Type: data.PARAMS}

	md := []*data.Attribute{attr1, attr2, attr3, attr4}
	outScope := data.NewFixedScope(md)

	objVal, _ := data.CoerceToObject("{\"key1\":5}")
	outScope.SetAttrValue("Obj", objVal)

	objVal, _ = data.CoerceToObject("{\"key1\":6}")
	outScope.SetAttrValue("Obj2", objVal)

	arrVal, _ := data.CoerceToArray("[1,6,3]")
	outScope.SetAttrValue("Array", arrVal)

	arrVal, _ = data.CoerceToArray("[7,8,9]")
	outScope.SetAttrValue("Array2", arrVal)

	paramVal, _ := data.CoerceToParams("{\"param1\":\"val\"}")
	outScope.SetAttrValue("Params", paramVal)

	paramVal, _ = data.CoerceToParams("{\"param1\":\"val2\"}")
	outScope.SetAttrValue("Params2", paramVal)


	err := mapper.Apply(nil, outScope)
	assert.Nil(t, err)

	expr := NewLookupExpr("${Obj}.key")
	newVal, err := expr.Eval(outScope)
	assert.Nil(t, err)
	assert.Equal(t, 2, newVal)

	expr = NewLookupExpr("${Array}[2]")
	newVal, err = expr.Eval(outScope)
	assert.Nil(t, err)
	assert.Equal(t, 3, newVal)

	expr = NewLookupExpr("${Params}.paramKey")
	newVal, err = expr.Eval(outScope)
	assert.Nil(t, err)
	assert.Equal(t, "4", newVal)
}

func TestAssignMapper(t *testing.T) {

	factory := GetFactory()

	mapping1 := &data.MappingDef{Type: data.MtAssign, Value: "${SimpleI}", MapTo: "SimpleO"}
	mapping2 := &data.MappingDef{Type: data.MtAssign, Value: "${ObjI}.key", MapTo: "ObjO.key"}
	mapping3 := &data.MappingDef{Type: data.MtAssign, Value: "${ArrayI}[2]", MapTo: "ArrayO[2]"}
	mapping4 := &data.MappingDef{Type: data.MtAssign, Value: "${ParamsI}.paramKey", MapTo: "ParamsO.paramKey"}

	mappings := []*data.MappingDef{mapping1, mapping2, mapping3, mapping4}

	mapper := factory.NewMapper(&data.MapperDef{Mappings: mappings})

	attrI1 := &data.Attribute{Name: "SimpleI", Type: data.INTEGER}
	attrI2 := &data.Attribute{Name: "ObjI", Type: data.OBJECT}
	attrI3 := &data.Attribute{Name: "ArrayI", Type: data.ARRAY}
	attrI4 := &data.Attribute{Name: "ParamsI", Type: data.PARAMS}

	mdI := []*data.Attribute{attrI1, attrI2, attrI3, attrI4}
	inScope := data.NewFixedScope(mdI)

	attrO1 := &data.Attribute{Name: "SimpleO", Type: data.INTEGER}
	attrO2 := &data.Attribute{Name: "ObjO", Type: data.OBJECT}
	attrO3 := &data.Attribute{Name: "ArrayO", Type: data.ARRAY}
	attrO4 := &data.Attribute{Name: "ParamsO", Type: data.PARAMS}

	mdO := []*data.Attribute{attrO1, attrO2, attrO3, attrO4}
	outScope := data.NewFixedScope(mdO)

	inScope.SetAttrValue("SimpleI", 1)

	objVal, _ := data.CoerceToObject("{\"key\":1}")
	inScope.SetAttrValue("ObjI", objVal)

	arrVal, _ := data.CoerceToArray("[1,2,3]")
	inScope.SetAttrValue("ArrayI", arrVal)

	paramVal, _ := data.CoerceToParams("{\"paramKey\":\"val1\"}")
	inScope.SetAttrValue("ParamsI", paramVal)

	objVal, _ = data.CoerceToObject("{\"key1\":5}")
	outScope.SetAttrValue("ObjO", objVal)

	arrVal, _ = data.CoerceToArray("[4,5,6]")
	outScope.SetAttrValue("ArrayO", arrVal)

	paramVal, _ = data.CoerceToParams("{\"param1\":\"val\"}")
	outScope.SetAttrValue("ParamsO", paramVal)

	err := mapper.Apply(inScope, outScope)
	assert.Nil(t, err)

	expr := NewLookupExpr("${ObjO}.key")
	newVal, err := expr.Eval(outScope)
	assert.Nil(t, err)
	assert.Equal(t, 1.0, newVal)

	expr = NewLookupExpr("${ArrayO}[2]")
	newVal, err = expr.Eval(outScope)
	assert.Nil(t, err)
	assert.Equal(t, 3.0, newVal)

	expr = NewLookupExpr("${ParamsO}.paramKey")
	newVal, err = expr.Eval(outScope)
	assert.Nil(t, err)
	assert.Equal(t, "val1", newVal)
}
