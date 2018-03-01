package direction

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/json/field"

	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/expr"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/function"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/gocc/token"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/ref"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

var log = logger.GetLogger("expression-direction")

type Attribute interface{}

func NewDoubleQuoteStringLit(lit interface{}) (string, error) {
	str := strings.TrimSpace(string(lit.(*token.Token).Lit))
	log.Debugf("New double qutoes [%s]", str)

	if str != "" && len(str) > 0 {
		str = field.RemoveQuote(str)
	}
	//Eascap string
	if strings.Contains(str, "\\\"") {
		str = strings.Replace(str, "\\\"", "\"", -1)
	}

	log.Debugf("Final double qutoes [%s]", str)

	return str, nil
}

func NewSingleQuoteStringLit(lit interface{}) (string, error) {
	str := strings.TrimSpace(string(lit.(*token.Token).Lit))
	log.Debugf("New single qutoe [%s]", str)

	if str != "" && len(str) > 0 {
		str = field.RemoveQuote(str)
	}

	//Eascap string
	if strings.Contains(str, "\\'") {
		str = strings.Replace(str, "\\'", "'", -1)
	}

	log.Debugf("Final single qutoe [%s]", str)

	return str, nil
}

func NewIntLit(lit interface{}) (int64, error) {
	str := strings.TrimSpace(string(lit.(*token.Token).Lit))
	s, err := strconv.ParseInt(str, 10, 64)
	return s, err
}

func NewBool(lit interface{}) (bool, error) {
	s := strings.TrimSpace(string(lit.(*token.Token).Lit))
	b, err := strconv.ParseBool(s)
	return b, err
}

func NewMappingRef(lit interface{}) (interface{}, error) {
	s := strings.TrimSpace(string(lit.(*token.Token).Lit))
	log.Debugf("New mapping ref and value [%s]", s)
	if strings.HasPrefix(s, "$.") || strings.HasPrefix(s, "$$") {
		m := ref.NewArrayRef(s)
		return m, nil
	} else {
		m := ref.NewMappingRef(s)
		return m, nil
	}
}

func NewFunction(name Attribute, parameters Attribute) (interface{}, error) {
	log.Debugf("New function name type [%s] and parameter type [%s]:", reflect.TypeOf(name), reflect.TypeOf(parameters))

	f_func := &function.FunctionExp{}
	to := name.(*token.Token)
	f_func.Name = string(to.Lit)

	switch parameters.(type) {
	case *function.Parameter:
		f_func.Params = append(f_func.Params, parameters.(*function.Parameter))
	case []*function.Parameter:
		for _, p := range parameters.([]*function.Parameter) {
			if !p.IsEmtpy() {
				f_func.Params = append(f_func.Params, p)
			}
		}
	}
	return f_func, nil
}

func NewArgument(a Attribute) (interface{}, error) {
	log.Debugf("New Argument and type [%s]", reflect.TypeOf(a))
	param := &function.Parameter{}
	parameters := []*function.Parameter{}
	switch a.(type) {
	case *token.Token:
		param.Type = data.STRING
		param.Value = string(a.(*token.Token).Lit)
	case int64:
		param.Type = data.INTEGER
		param.Value = a
	case string:
		param.Type = data.STRING
		param.Value = a.(string)
	case bool:
		param.Type = data.BOOLEAN
		param.Value = a.(bool)
	case *function.FunctionExp:
		param.Type = data.FUNCTION
		param.Function = a.(*function.FunctionExp)
	case []*function.Parameter:
		for _, p := range a.([]*function.Parameter) {
			if !p.IsEmtpy() {
				parameters = append(parameters, p)
			}
		}
	case *ref.MappingRef:
		param.Type = data.REF
		param.Value = a
	case *ref.ArrayRef:
		param.Type = data.ARRAYREF
		param.Value = a
	case []interface{}:
		//TODO
		log.Debug("New Arguments type is []interface{}")
	case interface{}:
		//TODO
		log.Debugf("New Arguments type is interface{}, [%+v]", reflect.TypeOf(a))
	}
	parameters = append(parameters, param)
	return parameters, nil
}

func NewArguments(as ...Attribute) (interface{}, error) {
	log.Debugf("New Arguments and type [%s]", reflect.TypeOf(as))
	parameters := []*function.Parameter{}
	for _, a := range as {
		param := &function.Parameter{}
		switch a.(type) {
		case *token.Token:
			param.Type = data.STRING
			param.Value = string(a.(*token.Token).Lit)
		case int64:
			param.Type = data.INTEGER
			param.Value = a
		case string:
			param.Type = data.STRING
			param.Value = a.(string)
		case *function.FunctionExp:
			param.Type = data.FUNCTION
			param.Function = a.(*function.FunctionExp)
		case *ref.MappingRef:
			param.Type = data.REF
			param.Value = a
		case *ref.ArrayRef:
			param.Type = data.ARRAYREF
			param.Value = a
		case []*function.Parameter:
			for _, p := range a.([]*function.Parameter) {
				if !p.IsEmtpy() {
					parameters = append(parameters, p)
				}
			}
		case []interface{}:
			log.Debugf("New Arguments type is []interface{}")
		case interface{}:
			log.Debugf("New Arguments type is interface{} %+v", a)
		}
		parameters = append(parameters, param)
	}
	return parameters, nil
}

func NewExpressionField(a Attribute) (interface{}, error) {
	log.Debugf("New Expression field [%+v] and type [%s]", a, reflect.TypeOf(a))
	expression := getExpression(a)
	return expression, nil
}

func NewExpression(left Attribute, op Attribute, right Attribute) (interface{}, error) {
	log.Debugf("New Expression and operator [%s]", string(op.(*token.Token).Lit))

	expression := expr.NewWIExpression()
	operator, found := expr.ToOperator(strings.TrimSpace(string(op.(*token.Token).Lit)))
	if found {
		expression.Operator = operator
	} else {
		return nil, errors.New("Unsupport operator " + string(op.(*token.Token).Lit))
	}

	expression.Left = getExpression(left)
	expression.Right = getExpression(right)
	expression.Type = data.EXPRESSION
	return expression, nil
}

func getExpression(ex Attribute) *expr.Expression {
	expression := expr.NewWIExpression()
	switch ex.(type) {
	case int64:
		expression.Type = data.INTEGER
		expression.Value = ex.(int64)
	case string:
		expression.Type = data.STRING
		expression.Value = ex.(string)
	case bool:
		expression.Type = data.BOOLEAN
		expression.Value = ex.(bool)
	case ref.MappingRef:
		expression.Type = data.REF
		ref := ex.(ref.MappingRef)
		expression.Value = ref
	case *ref.MappingRef:
		expression.Type = data.REF
		expression.Value = ex.(*ref.MappingRef).GetRef()
	case *ref.ArrayRef:
		expression.Type = data.ARRAYREF
		expression.Value = ex.(*ref.ArrayRef).GetRef()
	case *function.FunctionExp:
		expression.Type = data.FUNCTION
		expression.Value = ex.(*function.FunctionExp)
	case *expr.Expression:
		expression = ex.(*expr.Expression)
	default:
		log.Errorf("Unknow expression type [%s]", ex)
	}
	return expression
}

func NewTernaryExpression(first Attribute, second Attribute, third Attribute) (Attribute, error) {
	log.Debugf("first [%+v] and type [%s]", first, reflect.TypeOf(first))
	log.Debugf("second [%+v] and type [%s]", second, reflect.TypeOf(second))
	log.Debugf("third [%+v] and type [%s]", third, reflect.TypeOf(third))
	ternaryExp := &expr.TernaryExpressio{First: first, Second: second, Third: third}
	return ternaryExp, nil
}
