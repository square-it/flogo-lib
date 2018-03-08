package equals

import (
	"fmt"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression"
	"github.com/stretchr/testify/assert"
)

var s = &Equals{}

func TestStaticFunc_Starts_with(t *testing.T) {
	final1 := s.Eval("TIBCO FLOGO", "TIBCO", true)
	fmt.Println(final1)
	assert.Equal(t, false, final1)

	final2 := s.Eval("TIBCO", "tibco", true)
	fmt.Println(final2)
	assert.Equal(t, true, final2)

}

func TestExpression(t *testing.T) {
	fun, err := expression.NewFunctionExpression(`string.equals("TIBCO FLOGO", "TIBCO FLOGO", false)`).GetFunction()
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval()
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, true, v[0])
}

func TestExpressionIgnoreCase(t *testing.T) {
	fun, err := expression.NewFunctionExpression(`string.equals("TIBCO flogo", "TIBCO FLOGO", true)`).GetFunction()
	assert.Nil(t, err)
	assert.NotNil(t, fun)
	v, err := fun.Eval()
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, true, v[0])
}
