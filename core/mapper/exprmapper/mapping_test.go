package exprmapper
//
//import (
//	"encoding/json"
//	"os"
//	"strings"
//	"testing"
//
//	_ "git.tibco.com/git/product/ipaas/wi-contrib.git/function/number/len"
//	_ "git.tibco.com/git/product/ipaas/wi-contrib.git/function/string/concat"
//
//	"github.com/Sirupsen/logrus"
//	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
//	"github.com/TIBCOSoftware/flogo-lib/core/data"
//
//	"runtime/debug"
//
//	"github.com/TIBCOSoftware/flogo-lib/core/mapper/expression/expression/function"
//	"github.com/stretchr/testify/assert"
//)
//
//func TestDoMapping(t *testing.T) {
//
//	inputScope := TestScope{attrs: make(map[string]*data.Attribute)}
//
//	outputScope := TestScope{attrs: make(map[string]*data.Attribute)}
//
//	mapping := &data.MappingDef{
//		Type:  data.MtExpression,
//		Value: "$A3.Account.ID",
//		MapTo: "Account.ID2"}
//
//	mappings := []*data.MappingDef{}
//	mappings = append(mappings, mapping)
//
//	WIMapper := NewWIMapper(&definition.MapperDef{Mappings: mappings})
//	WIMapper.Apply(inputScope, outputScope)
//
//	attr, exist := outputScope.GetAttr("Account")
//	assert.Equal(t, true, exist)
//	assert.NotEmpty(t, attr.Value.(*data.ComplexObject).Value)
//	v, _ := json.Marshal(attr)
//	logrus.Infof("Final destionation field %s", string(v))
//
//}
//
//func TestDoMapping1Level(t *testing.T) {
//
//	inputScope := TestScope{attrs: make(map[string]*data.Attribute)}
//
//	outputScope := TestScope{attrs: make(map[string]*data.Attribute)}
//
//	mapping := &data.MappingDef{
//		Type:  data.MtExpression,
//		Value: "$A3.Account.Address.Streat",
//		MapTo: "Account.Address.State"}
//
//	mappings := []*data.MappingDef{}
//	mappings = append(mappings, mapping)
//
//	WIMapper := NewWIMapper(&definition.MapperDef{Mappings: mappings})
//	WIMapper.Apply(inputScope, outputScope)
//
//	attr, exist := outputScope.GetAttr("Account")
//	assert.Equal(t, true, exist)
//	assert.NotEmpty(t, attr.Value.(*data.ComplexObject).Value)
//	v, _ := json.Marshal(attr)
//	logrus.Infof("Final destionation field %s", string(v))
//
//}
//
//func TestDoMapping2Level(t *testing.T) {
//
//	inputScope := TestScope{attrs: make(map[string]*data.Attribute)}
//
//	outputScope := TestScope{attrs: make(map[string]*data.Attribute)}
//
//	mapping := &data.MappingDef{
//		Type:  data.MtExpression,
//		Value: "$A3.Account.Address.City.Park.Name",
//		MapTo: "Account.Address.State"}
//
//	mappings := []*data.MappingDef{}
//	mappings = append(mappings, mapping)
//
//	WIMapper := NewWIMapper(&definition.MapperDef{Mappings: mappings})
//	WIMapper.Apply(inputScope, outputScope)
//
//	attr, exist := outputScope.GetAttr("Account")
//	assert.Equal(t, true, exist)
//	assert.NotEmpty(t, attr.Value.(*data.ComplexObject).Value)
//	v, _ := json.Marshal(attr.Value.(*data.ComplexObject).Value)
//	logrus.Infof("Final destionation field %s", string(v))
//}
//
//func TestDoMappingFunction(t *testing.T) {
//
//	inputScope := TestScope{attrs: make(map[string]*data.Attribute)}
//
//	outputScope := TestScope{attrs: make(map[string]*data.Attribute)}
//
//	mapping := &data.MappingDef{
//		Type:  data.MtExpression,
//		Value: `string.concat($A3.Account.Address.City.Park.Name,"function string")`,
//		MapTo: "Account.Address.State"}
//
//	mappings := []*data.MappingDef{}
//	mappings = append(mappings, mapping)
//
//	WIMapper := NewWIMapper(&definition.MapperDef{Mappings: mappings})
//	WIMapper.Apply(inputScope, outputScope)
//
//	attr, exist := outputScope.GetAttr("Account")
//	assert.Equal(t, true, exist)
//	assert.NotEmpty(t, attr.Value.(*data.ComplexObject).Value)
//	v, _ := json.Marshal(attr.Value.(*data.ComplexObject).Value)
//	logrus.Infof("Final destionation field %s", string(v))
//
//}
//
//func TestDoMappingFunctionError(t *testing.T) {
//	defer func() {
//		if r := recover(); r != nil {
//			assert.NotNil(t, r)
//		}
//	}()
//	inputScope := TestScope{attrs: make(map[string]*data.Attribute)}
//
//	outputScope := TestScope{attrs: make(map[string]*data.Attribute)}
//
//	mapping := &data.MappingDef{
//		Type:  data.MtExpression,
//		Value: `string.concat("sss",123)`,
//		MapTo: "Account.Address.State"}
//
//	mappings := []*data.MappingDef{}
//	mappings = append(mappings, mapping)
//
//	WIMapper := NewWIMapper(&definition.MapperDef{Mappings: mappings})
//	WIMapper.Apply(inputScope, outputScope)
//
//}
//
//func TestDoMappingWithNull(t *testing.T) {
//
//	defer func() {
//		if r := recover(); r != nil {
//			log.Errorf("Apply mapping error %+v", r)
//			log.Errorf("StackTrace: %s", debug.Stack())
//		}
//	}()
//	inputScope := TestScope{attrs: make(map[string]*data.Attribute)}
//
//	outputScope := TestScope{attrs: make(map[string]*data.Attribute)}
//
//	mapping := &data.MappingDef{
//		Type:  data.MtExpression,
//		Value: `string.concat($A3.Account.ID2, "lixingwang")`,
//		MapTo: "$A4.Account.ID2"}
//
//	mappings := []*data.MappingDef{}
//	mappings = append(mappings, mapping)
//
//	WIMapper := NewWIMapper(&definition.MapperDef{Mappings: mappings})
//	WIMapper.Apply(inputScope, outputScope)
//
//	attr, exist := outputScope.GetAttr("Account")
//	assert.Equal(t, true, exist)
//	assert.NotEmpty(t, attr.Value.(*data.ComplexObject).Value)
//	v, _ := json.Marshal(attr.Value.(*data.ComplexObject).Value)
//	logrus.Infof("Final destionation field %s", string(v))
//
//}
//
//func TestDoMappingWithNullLastArgument(t *testing.T) {
//
//	defer func() {
//		if r := recover(); r != nil {
//			log.Errorf("Apply mapping error %+v", r)
//			log.Errorf("StackTrace: %s", debug.Stack())
//		}
//	}()
//	inputScope := TestScope{attrs: make(map[string]*data.Attribute)}
//
//	outputScope := TestScope{attrs: make(map[string]*data.Attribute)}
//
//	mapping := &data.MappingDef{
//		Type:  data.MtExpression,
//		Value: `string.test("ddddd", 123, $A3.Account.ID2)`,
//		MapTo: "$A4.Account.ID2"}
//
//	mappings := []*data.MappingDef{}
//	mappings = append(mappings, mapping)
//
//	WIMapper := NewWIMapper(&definition.MapperDef{Mappings: mappings})
//	WIMapper.Apply(inputScope, outputScope)
//
//	attr, exist := outputScope.GetAttr("Account")
//	assert.Equal(t, true, exist)
//	assert.NotEmpty(t, attr.Value.(*data.ComplexObject).Value)
//	v, _ := json.Marshal(attr.Value.(*data.ComplexObject).Value)
//	logrus.Infof("Final destionation field %s", string(v))
//
//}
//
//func TestDoMappingExpression(t *testing.T) {
//
//	defer func() {
//		if r := recover(); r != nil {
//			log.Errorf("Apply mapping error %+v", r)
//			log.Errorf("StackTrace: %s", debug.Stack())
//		}
//	}()
//	inputScope := TestScope{attrs: make(map[string]*data.Attribute)}
//
//	outputScope := TestScope{attrs: make(map[string]*data.Attribute)}
//
//	mapping := &data.MappingDef{
//		Type:  data.MtExpression,
//		Value: `number.len($A3.Account.Address.City.Park.Name)>2`,
//		MapTo: "$A4.Account.Address.IsCity"}
//
//	mappings := []*data.MappingDef{}
//	mappings = append(mappings, mapping)
//
//	WIMapper := NewWIMapper(&definition.MapperDef{Mappings: mappings})
//	WIMapper.Apply(inputScope, outputScope)
//
//	attr, exist := outputScope.GetAttr("Account")
//	assert.Equal(t, true, exist)
//	assert.NotEmpty(t, attr.Value.(*data.ComplexObject).Value)
//	v, _ := json.Marshal(attr.Value.(*data.ComplexObject).Value)
//	logrus.Infof("Final destionation field %s", string(v))
//
//}
//
//func TestDoMappingArray(t *testing.T) {
//
//	//TODO handle array
//
//	//t.Fatal("Need Hanle array")
//
//	inputScope := TestScope{attrs: make(map[string]*data.Attribute)}
//
//	outputScope := TestScope{attrs: make(map[string]*data.Attribute)}
//
//	mapping := &data.MappingDef{
//		Type:  data.MtExpression,
//		Value: `$A3.Account.Emails`,
//		MapTo: "$A4.Account.EmailsSet"}
//
//	mappings := []*data.MappingDef{}
//	mappings = append(mappings, mapping)
//
//	WIMapper := NewWIMapper(&definition.MapperDef{Mappings: mappings})
//	WIMapper.Apply(inputScope, outputScope)
//
//	attr, exist := outputScope.GetAttr("Account")
//	assert.Equal(t, true, exist)
//	assert.NotEmpty(t, attr.Value.(*data.ComplexObject).Value)
//
//	v, _ := json.Marshal(attr.Value.(*data.ComplexObject).Value)
//	logrus.Infof("Final destionation field %s", string(v))
//
//}
//
//func TestSpecialFieldMapping(t *testing.T) {
//
//	//TODO handle array
//	os.Setenv("FLOGO_LOG_LEVEL", "debug")
//	//t.Fatal("Need Hanle array")
//
//	inputScope := TestScope{attrs: make(map[string]*data.Attribute)}
//
//	outputScope := TestScope{attrs: make(map[string]*data.Attribute)}
//
//	mapping := &data.MappingDef{
//		Type:  data.MtExpression,
//		Value: `$A3.Account.Maps3["dd*cc"]["y.x"][0][d.d]`,
//		MapTo: "$A4.Account.Maps"}
//
//	mappings := []*data.MappingDef{}
//	mappings = append(mappings, mapping)
//
//	WIMapper := NewWIMapper(&definition.MapperDef{Mappings: mappings})
//	WIMapper.Apply(inputScope, outputScope)
//
//	attr, exist := outputScope.GetAttr("Account")
//	assert.Equal(t, true, exist)
//	assert.NotEmpty(t, attr.Value.(*data.ComplexObject).Value)
//
//	v, _ := json.Marshal(attr.Value.(*data.ComplexObject).Value)
//	logrus.Infof("Final destionation field %s", string(v))
//
//}
//
//type TestScope struct {
//	attrs map[string]*data.Attribute
//}
//
//// GetAttr gets the specified attribute
//func (s TestScope) GetAttr(name string) (attr *data.Attribute, exists bool) {
//	if s.attrs != nil && len(s.attrs) > 0 && s.attrs[name] != nil {
//		return s.attrs[name], true
//	} else {
//		if strings.EqualFold(name, "test") {
//			complexObject := data.ComplexObject{
//				Metadata: "",
//				Value:    AccountStr}
//
//			s.attrs[name] = data.NewAttribute("Account", data.COMPLEX_OBJECT, &complexObject)
//			return s.attrs[name], true
//
//		} else {
//			complexObject := data.ComplexObject{
//				Metadata: "",
//				Value:    AccountStr}
//
//			s.attrs[name] = data.NewAttribute("Account", data.COMPLEX_OBJECT, &complexObject)
//			return s.attrs[name], true
//		}
//
//		return nil, false
//	}
//
//}
//
//func stringP(s string) *string {
//	return &s
//}
//
//// SetAttrValue sets the value of the specified attribute
//func (s TestScope) SetAttrValue(name string, value interface{}) error {
//	logrus.Infof("SetAttrValue Name %s and value %+v", name, value)
//	if s.attrs == nil {
//		s.attrs = make(map[string]*data.Attribute)
//	}
//
//	existingAttr, exists := s.GetAttr(name)
//
//	//todo: optimize, use existing attr
//	if exists {
//		attr := data.NewAttribute(name, existingAttr.Type, value)
//		s.attrs[name] = attr
//		return nil
//	} else {
//		logrus.Info("Not exist !")
//	}
//
//	return nil
//
//}
//
//var AccountStr = `{
//    "Address": {
//        "City": {
//            "InUS": true,
//            "Name": "Sugar Land",
//            "Park": {
//                "Location": "location",
//                "Name": "Name"
//            }
//        },
//        "IsCity": false,
//        "State": "Taxes",
//        "Streat": "311 wind st",
//        "ZipCode": "77477"
//    },
//    "Emails": [
//        "123@gmail.com",
//        "234@gmail.com"
//    ],
//    "EmailsSet": null,
//    "ID": "10001",
//	"ID2": null,
//	"Maps3": {
//        "bb.bb": [
//            {
//                "x.y": "10001"
//            }
//        ],
//        "cc#cc": [
//            {
//                "x.y": {
//                    "id": "1"
//                }
//            }
//        ],
//        "dd*cc": {
//            "x.y": {
//                "g%f": "hello"
//            },
//            "y.x":[
//                {"d.d":"123"}
//            ]
//        }
//    }
//}`
//
//type Test struct {
//}
//
//func init() {
//	function.Registry(&Test{})
//}
//
//func (s *Test) GetName() string {
//	return "test"
//}
//
//func (s *Test) GetCategory() string {
//	return "string"
//}
//
//func (s *Test) Eval(name string, age int64, names ...string) string {
//	value := name + strings.Join(names, ":")
//	log.Infof("Resut %s", value)
//	return value
//}
//
//func TestRemoveMapToPrefic(t *testing.T) {
//	mapping := `{
//    "fields": [
//        {
//            "from": "string.concat(\"hahahahahah\",$.$ref)",
//            "to": "$INPUT.$$$ref",
//            "type": "primitive"
//        },
//        {
//            "from": "$._name",
//            "to": "$INPUT._name",
//            "type": "primitive"
//        },
//        {
//            "from": "$.last",
//            "to": "$INPUT.$$last",
//            "type": "primitive"
//        },
//        {
//            "fields": [
//                {
//                    "from": "string.concat(\"hahahahahah\",$.$ref)",
//                    "to": "$INPUT.$$$ref",
//                    "type": "primitive"
//                },
//                {
//                    "from": "$._name",
//                    "to": "$INPUT._name",
//                    "type": "primitive"
//                },
//                {
//                    "from": "$.last",
//                    "to": "$INPUT.$$last",
//                    "type": "primitive"
//                },
//                {
//                    "fields": [
//                        {
//                            "from": "string.concat(\"hahahahahah\",$.$ref)",
//                            "to": "$INPUT.$$$ref",
//                            "type": "primitive"
//                        },
//                        {
//                            "from": "$._name",
//                            "to": "$INPUT.$$_name",
//                            "type": "primitive"
//                        },
//                        {
//                            "from": "$.last",
//                            "to": "$INPUT.$$last",
//                            "type": "primitive"
//                        }
//                    ],
//                    "from": "$INPUT$.resources2",
//                    "to": "data.resources",
//                    "type": "foreach"
//                }
//            ],
//            "from": "$TriggerData.body.resources",
//            "to": "$INPUT.$$[\"data\"][\"resources\"]",
//            "type": "foreach"
//        }
//    ],
//    "from": "$TriggerData.body.resources",
//    "to": "data.resources",
//    "type": "foreach"
//}`
//
//	a, err := ParseArrayMapping(mapping)
//	assert.Nil(t, err)
//	a.RemovePrefixForMapTo()
//	assert.Equal(t, "$$$ref", a.Fields[3].Fields[3].Fields[0].To)
//	assert.Equal(t, "$$_name", a.Fields[3].Fields[3].Fields[1].To)
//	assert.Equal(t, "$$last", a.Fields[3].Fields[3].Fields[2].To)
//
//}
