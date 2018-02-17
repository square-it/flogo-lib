package exprmapper
//
//import (
//	"strings"
//	"testing"
//
//	"fmt"
//
//	_ "git.tibco.com/git/product/ipaas/wi-contrib.git/function/number/Int64"
//	_ "git.tibco.com/git/product/ipaas/wi-contrib.git/function/string/concat"
//	_ "git.tibco.com/git/product/ipaas/wi-contrib.git/function/string/stringlength"
//	_ "git.tibco.com/git/product/ipaas/wi-contrib.git/function/string/substring"
//	_ "git.tibco.com/git/product/ipaas/wi-contrib.git/function/string/substringafter"
//	_ "git.tibco.com/git/product/ipaas/wi-contrib.git/function/string/substringbefore"
//	_ "git.tibco.com/git/product/ipaas/wi-contrib.git/function/string/tostring"
//	"github.com/TIBCOSoftware/flogo-lib/core/data"
//	"github.com/stretchr/testify/assert"
//)
//
//var maping = `{
//    "fields": [
//        {
//            "from": "string.concat($.name,\"_suffix\")",
//            "to": "$.name",
//            "type": "primitive"
//        },
//        {
//            "from": "string.concat(\"countrycode\",$.phone)",
//            "to": "$.phone",
//            "type": "primitive"
//        },
//        {
//            "fields": [
//                {
//                    "from": "$.city",
//                    "to": "$.city",
//                    "type": "primitive"
//                },
//                {
//                    "from": "$.zipcode",
//                    "to": "$.zipcode",
//                    "type": "primitive"
//                }
//            ],
//            "from": "$.Address",
//            "to": "$.Address",
//            "type": "foreach"
//        }
//    ],
//    "from": "$T.Persons",
//    "to": "$.students",
//    "type": "foreach"
//}`
//
//
//type TestScopeArray struct {
//	attrs map[string]*data.Attribute
//}
//
//// GetAttr gets the specified attribute
//func (s TestScopeArray) GetAttr(name string) (attr *data.Attribute, exists bool) {
//	if s.attrs != nil && len(s.attrs) > 0 && s.attrs[name] != nil {
//		return s.attrs[name], true
//	} else {
//		if strings.EqualFold(name, "{A3.Account}") {
//			complexObject := data.ComplexObject{
//				Metadata: "",
//				Value:    AccountStrSource}
//
//			s.attrs[name] = data.NewAttribute("Account", data.COMPLEX_OBJECT, &complexObject)
//			return s.attrs[name], true
//
//		} else {
//			complexObject := data.ComplexObject{
//				Metadata: "",
//				Value:    AccountDestination}
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
//// SetAttrValue sets the value of the specified attribute
//func (s TestScopeArray) SetAttrValue(name string, value interface{}) error {
//	arraylog.Infof("SetAttrValue Name %s and value %+v", name, value)
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
//		arraylog.Info("Not exist !")
//	}
//
//	return nil
//
//}
//
//var AccountStrSource = `{
//    "Address": {
//        "City": [
//            {
//                "InUS": true,
//                "Name": "Sugar Land",
//                "Park": {
//                    "Location": "sugar land street",
//                    "Name": "sugarland park"
//                },
//                "TestArray":[
//                {"name1":"name1-1"},
//                {"name1":"name2-2"},
//                {"name1":"name3-3"}
//                ]
//            },
//            {
//                "InUS": true,
//                "Name": "Stafford",
//                "Park": {
//                    "Location": "Stafford street",
//                    "Name": "Stafford park"
//                },
//                "TestArray":[
//                {"name1":"name1-2"},
//                {"name1":"name2-2"},
//                {"name1":"name3-3"}
//                ]
//            }
//        ],
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
//    "ID2": null
//}`
//
//var AccountDestination = `{}`
//
//func TestToArrayMapping(t *testing.T) {
//	v, err := ParseArrayMapping(maping)
//	assert.Nil(t, err)
//	assert.NotNil(t, v)
//}
//
//func TestArrayMapping1Level(t *testing.T) {
//
//	arraymaping := `
//	{
//	   "fields":[
//	      {
//		 "from":"$.Name",
//		 "to":"$.Name",
//		 "type":"primitive"
//	      },
//	      {
//		 "from":"$.InUS",
//		 "to":"$.InUS",
//		 "type":"primitive"
//	      }
//	   ],
//	   "from":"$A3.Account.Address.City",
//	   "to":"Account.Address.City",
//	   "type":"foreach"
//	}
//	`
//
//	inputScope := TestScopeArray{attrs: make(map[string]*data.Attribute)}
//	outputScope := TestScopeArray{attrs: make(map[string]*data.Attribute)}
//	arrayMapping, err := ParseArrayMapping(arraymaping)
//	assert.Nil(t, err)
//	err = arrayMapping.DoArrayMapping(inputScope, outputScope)
//	if err != nil {
//		arraylog.Errorf("Do array mapping error %s", err)
//	}
//	assert.Nil(t, err)
//
//}
//
//func TestArrayMappingFunctionOnLeaf(t *testing.T) {
//
//	arraymaping := `
//	{
//	   "fields":[
//	      {
//		 "from":"string.concat($.Name,\"lixingwang\")",
//		 "to":"$.Name",
//		 "type":"primitive"
//	      },
//	      {
//		 "from":"$.InUS",
//		 "to":"$.InUS",
//		 "type":"primitive"
//	      }
//	   ],
//	   "from":"$A3.Account.Address.City",
//	   "to":"Account.Address.City",
//	   "type":"foreach"
//	}
//	`
//
//	inputScope := TestScopeArray{attrs: make(map[string]*data.Attribute)}
//	outputScope := TestScopeArray{attrs: make(map[string]*data.Attribute)}
//	arrayMapping, err := ParseArrayMapping(arraymaping)
//	assert.Nil(t, err)
//	err = arrayMapping.DoArrayMapping(inputScope, outputScope)
//	arraylog.Errorf("Do array mapping error %s", err)
//	assert.Nil(t, err)
//
//}
//
//func TestArrayMappingFunctionLiteral(t *testing.T) {
//
//	arraymaping := `
//	{
//	   "fields":[
//	      {
//		 "from":"wangzai",
//		 "to":"$.Name",
//		 "type":"primitive"
//	      },
//	      {
//		 "from":"$.InUS",
//		 "to":"$.InUS",
//		 "type":"primitive"
//	      }
//	   ],
//	   "from":"$A3.Account.Address.City",
//	   "to":"Account.Address.City",
//	   "type":"foreach"
//	}
//	`
//
//	inputScope := TestScopeArray{attrs: make(map[string]*data.Attribute)}
//	outputScope := TestScopeArray{attrs: make(map[string]*data.Attribute)}
//	arrayMapping, err := ParseArrayMapping(arraymaping)
//	assert.Nil(t, err)
//	err = arrayMapping.DoArrayMapping(inputScope, outputScope)
//	attr, _ := outputScope.GetAttr("Account")
//	assert.NotNil(t, attr)
//	v, _ := attr.MarshalJSON()
//	fmt.Println("dddddddd", string(v))
//	assert.Nil(t, err)
//
//}
//
//func TestArrayMappingFunction2Level(t *testing.T) {
//
//	arraymaping := `
//	{
//	    "fields": [
//		{
//		    "from": "string.concat($.Name,\"lixingwang\")",
//		    "to": "$.Name",
//		    "type": "primitive"
//		},
//		{
//		    "from": "$.InUS",
//		    "to": "$.InUS",
//		    "type": "primitive"
//		},
//		{
//		    "fields": [
//			{
//			    "from": "$.name1",
//			    "to": "$.name1",
//			    "type": "primitive"
//			}
//		    ],
//		    "from": "$.TestArray",
//		    "to": "$.TestArray",
//		    "type": "foreach"
//		}
//	    ],
//	    "from": "$A3.Account.Address.City",
//	    "to": "Account.Address.City",
//	    "type": "foreach"
//	}`
//
//	inputScope := TestScopeArray{attrs: make(map[string]*data.Attribute)}
//	outputScope := TestScopeArray{attrs: make(map[string]*data.Attribute)}
//	arrayMapping, err := ParseArrayMapping(arraymaping)
//	assert.Nil(t, err)
//	assert.Nil(t, arrayMapping.Validate())
//	err = arrayMapping.DoArrayMapping(inputScope, outputScope)
//	assert.Nil(t, err)
//
//}
//
//func TestArrayMappingNewARRAY(t *testing.T) {
//
//	arraymaping := `
//	{
//	    "fields": [
//		{
//		    "from": "string.concat($.Name,\"lixingwang\")",
//		    "to": "$.Name",
//		    "type": "primitive"
//		},
//		{
//		    "from": "no",
//		    "to": "$.InUS",
//		    "type": "primitive"
//		},
//		{
//		    "fields": [
//			{
//			    "from": "$.name1",
//			    "to": "$.name1",
//			    "type": "primitive"
//			}
//		    ],
//		    "from": "NEWARRAY",
//		    "to": "$.TestArray",
//		    "type": "foreach"
//		}
//	    ],
//	    "from": "NEWARRAY",
//	    "to": "Account.Address.City",
//	    "type": "foreach"
//	}`
//
//	arrayMapping, err := ParseArrayMapping(arraymaping)
//	assert.Nil(t, err)
//	err = arrayMapping.Validate()
//	assert.NotNil(t, err)
//	fmt.Println(err)
//
//}
//
//func TestArrayMappingNewARRAY2(t *testing.T) {
//
//	arraymaping := `
//	{
//	    "fields": [
//		{
//		    "from": "string.concat($.Name,\"lixingwang\")",
//		    "to": "$.Name",
//		    "type": "primitive"
//		},
//		{
//		    "from": "no",
//		    "to": "$.InUS",
//		    "type": "primitive"
//		},
//		{
//		    "fields": [
//			{
//			    "from": "$.name1",
//			    "to": "$.name1",
//			    "type": "primitive"
//			}
//		    ],
//		    "from": "NEWARRAY",
//		    "to": "$.TestArray",
//		    "type": "foreach"
//		}
//	    ],
//	   "from":"$A3.Account.Address.City",
//	    "to": "Account.Address.City",
//	    "type": "foreach"
//	}`
//
//	arrayMapping, err := ParseArrayMapping(arraymaping)
//	assert.Nil(t, err)
//	err = arrayMapping.Validate()
//	assert.NotNil(t, err)
//
//}
//
//func TestArrayMappingNewARRAY3(t *testing.T) {
//
//	arraymaping := `
//	{
//	    "fields": [
//		{
//		    "from": "string.concat($.Name,\"lixingwang\")",
//		    "to": "$.Name",
//		    "type": "primitive"
//		},
//		{
//		    "from": "no",
//		    "to": "$.InUS",
//		    "type": "primitive"
//		},
//		{
//		    "fields": [
//			{
//			    "from": "$.name1",
//			    "to": "$.name1",
//			    "type": "primitive"
//			},
//	        {
//			    "from": 1222,
//			    "to": "$.street",
//			    "type": "primitive"
//			}
//		    ],
//		    "from": "NEWARRAY",
//		    "to": "$.TestArray",
//		    "type": "foreach"
//		}
//	    ],
//	   "from":"$A3.Account.Address.City",
//	    "to": "Account.Address.City",
//	    "type": "foreach"
//	}`
//
//	arrayMapping, err := ParseArrayMapping(arraymaping)
//	assert.Nil(t, err)
//	err = arrayMapping.Validate()
//	assert.NotNil(t, err)
//
//}
