package expression

import (
	"fmt"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/expr"
)


func TestFuncExpression(t *testing.T) {
	ex, err := NewExpression(`string.concat("123","456")="123456"`).GetExpression()
	if err != nil {
		fail(t, err)
	}

	exprStr, err := ex.Serialization()
	if err != nil {
		fail(t, err)
	}

	exprSer := &expr.Expression{}
	exprSer, err = expr.DeSerialization(exprStr)
	if err != nil {
		fail(t, err)
	}

	v, err := exprSer.Eval()
	if err != nil {
		fail(t, err)
	}

	fmt.Println("Result:", v)
}

func TestFuncExpressionAnd(t *testing.T) {
	ex, err := NewExpression(`("dddddd" = "dddddd") & ("123" = "123")`).GetExpression()
	if err != nil {
		fail(t, err)
	}

	exprStr, err := ex.Serialization()
	if err != nil {
		fail(t, err)
	}

	exprSer := &expr.Expression{}
	exprSer, err = expr.DeSerialization(exprStr)
	if err != nil {
		fail(t, err)
	}

	v, err := exprSer.Eval()
	if err != nil {
		fail(t, err)
	}

	fmt.Println("Result:", v)
}

func fail(t *testing.T, err error) {
	panic(err)
	t.Fatal(err)
	t.Failed()

}
