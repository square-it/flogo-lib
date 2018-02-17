package expression
//
//import (
//	"bytes"
//	"fmt"
//	"testing"
//
//	"github.com/TIBCOSoftware/flogo-lib/core/mapper/expression/expression/function"
//	_ "git.tibco.com/git/product/ipaas/wi-contrib.git/function/datetime/currentdate"
//	_ "git.tibco.com/git/product/ipaas/wi-contrib.git/function/datetime/currentdatetime"
//	_ "git.tibco.com/git/product/ipaas/wi-contrib.git/function/datetime/currenttime"
//	_ "git.tibco.com/git/product/ipaas/wi-contrib.git/function/number/Int64"
//	_ "git.tibco.com/git/product/ipaas/wi-contrib.git/function/string/tostring"
//	_ "git.tibco.com/git/product/ipaas/wi-contrib.git/function/string/concat"
//	_ "git.tibco.com/git/product/ipaas/wi-contrib.git/function/string/split"
//	_ "git.tibco.com/git/product/ipaas/wi-contrib.git/function/string/startswith"
//	_ "git.tibco.com/git/product/ipaas/wi-contrib.git/function/string/stringlength"
//	_ "git.tibco.com/git/product/ipaas/wi-contrib.git/function/string/substring"
//	_ "git.tibco.com/git/product/ipaas/wi-contrib.git/function/string/substringafter"
//	"github.com/stretchr/testify/assert"
//)
//
////
////func TestCallConcat(t *testing.T) {
////
////	f := &function.Function{Name: "concat", Params: []*function.Parameter{{Type: witype.STRING, Value: "li"}, {Type: witype.STRING, Value: "xingwang"}}}
////
////	v, err := f.Eval()
////	if err != nil {
////		t.Fatal(err)
////		t.Failed()
////	}
////	fmt.Println("Result:", v)
////}
//
//func TestCurrentDate(t *testing.T) {
//	v, err := NewFunctionExpression(`datetime.currentDate()`).Eval()
//	if err != nil {
//		t.Fatal(err)
//		t.Failed()
//	}
//	assert.NotNil(t, v)
//	fmt.Println("Result:", v)
//}
//
//func TestCurrentDatetimeAndTime(t *testing.T) {
//	v, err := NewFunctionExpression(`datetime.currentDatetime()`).Eval()
//	if err != nil {
//		t.Fatal(err)
//		t.Failed()
//	}
//	assert.NotNil(t, v)
//	fmt.Println("Result:", v)
//
//	v, err = NewFunctionExpression(`datetime.currentTime()`).Eval()
//	if err != nil {
//		t.Fatal(err)
//		t.Failed()
//	}
//	assert.NotNil(t, v)
//	fmt.Println("Result:", v)
//}
//
//func TestFunctionConcatWithSpace(t *testing.T) {
//
//	v, err := NewFunctionExpression(`string.concat("This", "is",string.concat("my","first"),"gocc",string.concat("lexer","and","parser"),string.concat("go","program","!!!"))`).Eval()
//	if err != nil {
//		t.Fatal(err)
//		t.Failed()
//	}
//
//	assert.Equal(t, "Thisismyfirstgocclexerandparsergoprogram!!!", v[0].(string))
//	fmt.Println("Result:", v)
//}
//
//func TestFunctionConcatWithMultiSpace(t *testing.T) {
//
//	v, err := NewFunctionExpression(`string.concat("This",   " is" , " WI")`).Eval()
//	if err != nil {
//		t.Fatal(err)
//		t.Failed()
//	}
//
//	assert.Equal(t, "This is WI", v[0].(string))
//	fmt.Println("Result:", v)
//}
//func TestFunctionConcat(t *testing.T) {
//
//	v, err := NewFunctionExpression(`string.concat("This","is",string.concat("my","first"),"gocc",string.concat("lexer","and","parser"),string.concat("go","program","!!!"))`).Eval()
//	if err != nil {
//		t.Fatal(err)
//		t.Failed()
//	}
//
//	assert.Equal(t, "Thisismyfirstgocclexerandparsergoprogram!!!", v[0].(string))
//	fmt.Println("Result:", v)
//}
//
//func TestFunctionSubstring(t *testing.T) {
//	v, err := NewFunctionExpression(`string.substring("lixingwang",2,5)`).Eval()
//	if err != nil {
//		t.Fatal(err)
//		t.Failed()
//	}
//	assert.Equal(t, "xingw", v[0].(string))
//
//	fmt.Println("Result:", v)
//}
//
//func TestFunctionLength(t *testing.T) {
//	v, err := NewFunctionExpression(`string.length("lixingwang")`).Eval()
//	if err != nil {
//		t.Fatal(err)
//		t.Failed()
//	}
//	assert.Equal(t, int64(10), v[0].(int64))
//
//	fmt.Println("Result:", v[0])
//}
//
//func TestFunctionCombine(t *testing.T) {
//	v, err := NewFunctionExpression(`string.concat("Beijing",string.tostring(string.length("lixingwang")))`).Eval()
//	if err != nil {
//		t.Fatal(err)
//		t.Failed()
//	}
//	assert.Equal(t, "Beijing10", v[0].(string))
//
//	fmt.Println("Result:", v[0])
//}
//
//func TestFunctionCombine2(t *testing.T) {
//	v, err := NewFunctionExpression(`string.concat("Beijing",string.tostring(number.int64(string.length("lixingwang"))))`).Eval()
//	if err != nil {
//		t.Fatal(err)
//		t.Failed()
//	}
//	assert.Equal(t, "Beijing10", v[0].(string))
//
//	fmt.Println("Result:", v[0])
//}
//
//func TestFunctionError(t *testing.T) {
//	v, err := NewFunctionExpression(`string.concat("Beijing",string.tostring(number.int64("2017")))`).Eval()
//	if err != nil {
//		assert.NotNil(t, err)
//		fmt.Println("Result", v)
//	} else {
//		t.Failed()
//	}
//}
//
//func TestFunctionWithRefMapping(t *testing.T) {
//	v, err := NewFunctionExpression(`string.concat($A3.query.result,"data")`).Eval()
//	if err != nil {
//		t.Fatal(err)
//		t.Failed()
//	}
//	assert.Equal(t, "$A3.query.resultdata", v[0].(string))
//
//	fmt.Println("Result:", v[0])
//}
//
//func TestFunctionWithRefMapping2(t *testing.T) {
//	v, err := NewFunctionExpression(`string.concat($A2.message,"lixingwang")`).Eval()
//	if err != nil {
//		t.Fatal(err)
//		t.Failed()
//	}
//	assert.Equal(t, "$A2.messagelixingwang", v[0].(string))
//
//	fmt.Println("Result:", v[0])
//}
//
//func TestFunctionWithTag(t *testing.T) {
//	v, err := NewFunctionExpression(`wi.concat($A2.message,"lixingwang")`).Eval()
//	if err != nil {
//		t.Fatal(err)
//		t.Failed()
//	}
//	assert.Equal(t, "$A2.messagelixingwang", v[0].(string))
//
//	fmt.Println("Result:", v[0])
//}
//
//func TestFunctionWithSpaceInRef(t *testing.T) {
//	v, err := NewFunctionExpression(`string.concat($Marketo Get Lead by Id.output.result[0].firstName,$Marketo Get Lead by Id.output.result[0].lastName)`).Eval()
//	if err != nil {
//		t.Fatal(err)
//		t.Failed()
//	}
//	assert.Equal(t, "$Marketo Get Lead by Id.output.result[0].firstName$Marketo Get Lead by Id.output.result[0].lastName", v[0].(string))
//
//	fmt.Println("Result:", v[0])
//}
//
//func TestFunctionWithPackage(t *testing.T) {
//	v, err := NewFunctionExpression(`string.concat($A2.message,"lixingwang")`).Eval()
//	if err != nil {
//		t.Fatal(err)
//		t.Failed()
//	}
//	assert.Equal(t, "$A2.messagelixingwang", v[0].(string))
//
//	fmt.Println("Result:", v[0])
//}
//
//func TestDashFunction(t *testing.T) {
//	v, err := NewFunctionExpression(`string.startsWith("TIBCOBW","TIBCO")`).Eval()
//	if err != nil {
//		t.Fatal(err)
//		t.Failed()
//	}
//	assert.Equal(t, true, v[0].(bool))
//
//	fmt.Println("Result:", v[0])
//}
//
//type Concat struct {
//}
//
//func init() {
//	function.Registry(&Concat{})
//}
//
//func (s *Concat) GetName() string {
//	return "concat"
//}
//
//func (s *Concat) GetCategory() string {
//	return "wi"
//}
//
//func (s *Concat) Eval(strs ...string) string {
//	log.Debugf("Start wi:concat function with parameters %s", strs)
//	var buffer bytes.Buffer
//
//	for _, v := range strs {
//		buffer.WriteString(v)
//	}
//	log.Debugf("Done wi:concat function with result %s", buffer.String())
//	return buffer.String()
//}
//
//type PConcat struct {
//}
//
//func init() {
//	function.Registry(&PConcat{})
//}
//
//func (s *PConcat) GetName() string {
//	return "string.concat"
//}
//
//func (s *PConcat) GetCategory() string {
//	return ""
//}
//
//func (s *PConcat) Eval(strs ...string) string {
//	log.Debugf("Start wi:concat function with parameters %s", strs)
//	var buffer bytes.Buffer
//
//	for _, v := range strs {
//		buffer.WriteString(v)
//	}
//	log.Debugf("Done wi:concat function with result %s", buffer.String())
//	return buffer.String()
//}
//
//type PanicFunc struct {
//}
//
//func init() {
//	function.Registry(&PanicFunc{})
//}
//
//func (s *PanicFunc) GetName() string {
//	return "panic"
//}
//
//func (s *PanicFunc) GetCategory() string {
//	return "panic"
//}
//
//func (s *PanicFunc) Eval() string {
//	panic("Panic happend")
//	return "panic"
//}
//
//func TestPanictFunction(t *testing.T) {
//	v, err := NewFunctionExpression(`panic.panic()`).Eval()
//	assert.NotNil(t, err)
//	assert.Nil(t, v)
//}
//
//func TestSplittFunction(t *testing.T) {
//	v, err := NewFunctionExpression(`string.split("hello,world",",")`).Eval()
//	assert.NotNil(t, v)
//	assert.Nil(t, err)
//	log.Info(v)
//}
//
//func TestNumberLenFunction(t *testing.T) {
//	v, err := NewFunctionExpression(`string.length("hello,world")`).Eval()
//	log.Info(v)
//	assert.NotNil(t, v)
//	assert.Nil(t, err)
//
//}
