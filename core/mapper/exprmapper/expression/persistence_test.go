package expression
//
//import (
//	"fmt"
//	"testing"
//
//	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/expr"
//	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/function"
//)
//
//func TestFunctionSubstringSer(t *testing.T) {
//	newFunction, err := NewFunctionExpression(`string.substring("lixingwang",2,5)`).GetFunction()
//	if err != nil {
//		fail(t, err)
//	}
//
//	str, err := newFunction.Serialization()
//	if err != nil {
//		fail(t, err)
//	}
//	funcUnseri := &function.FunctionExp{}
//	funcUnseri, err = function.DeSerialization(str)
//	v, err := funcUnseri.Eval()
//	if err != nil {
//		fail(t, err)
//	}
//
//	fmt.Println("Result:", v)
//}
//
//func TestFuncNestConcat(t *testing.T) {
//	newFunction, err := NewFunctionExpression(`string.concat("This","is",string.concat("my","first"),"gocc",string.concat("lexer","and","parser"),string.concat("go","program","!!!"))`).GetFunction()
//	if err != nil {
//		fail(t, err)
//	}
//
//	str, err := newFunction.Serialization()
//	if err != nil {
//		fail(t, err)
//	}
//	funcUnseri := &function.FunctionExp{}
//	funcUnseri, err = function.DeSerialization(str)
//	v, err := funcUnseri.Eval()
//	if err != nil {
//		fail(t, err)
//	}
//
//	fmt.Println("Result:", v)
//}
//
//func TestFuncExpression(t *testing.T) {
//	ex, err := NewExpression(`string.concat("123","456")="123456"`).GetExpression()
//	if err != nil {
//		fail(t, err)
//	}
//
//	exprStr, err := ex.Serialization()
//	if err != nil {
//		fail(t, err)
//	}
//
//	exprSer := &expr.Expression{}
//	exprSer, err = expr.DeSerialization(exprStr)
//	if err != nil {
//		fail(t, err)
//	}
//
//	v, err := exprSer.Eval()
//	if err != nil {
//		fail(t, err)
//	}
//
//	fmt.Println("Result:", v)
//}
//
//func TestFuncExpressionAnd(t *testing.T) {
//	ex, err := NewExpression(`("dddddd" = "dddddd") & ("123" = "123")`).GetExpression()
//	if err != nil {
//		fail(t, err)
//	}
//
//	exprStr, err := ex.Serialization()
//	if err != nil {
//		fail(t, err)
//	}
//
//	exprSer := &expr.Expression{}
//	exprSer, err = expr.DeSerialization(exprStr)
//	if err != nil {
//		fail(t, err)
//	}
//
//	v, err := exprSer.Eval()
//	if err != nil {
//		fail(t, err)
//	}
//
//	fmt.Println("Result:", v)
//}
//
//func fail(t *testing.T, err error) {
//	panic(err)
//	t.Fatal(err)
//	t.Failed()
//
//}
