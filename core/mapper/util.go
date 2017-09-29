package mapper

import (
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

func NewLookupExpr(exprRep string) data.Expr {
	return &lookupExpr{rep: exprRep}
}

type lookupExpr struct {
	rep string
}

func (e *lookupExpr) Eval(scope data.Scope) (interface{}, error) {

	resType, attrName, path, err := data.GetResolutionInfo(e.rep)

	if err != nil {
		return nil, err
	}

	var value interface{}
	var exists bool

	switch resType {
	case data.RES_PROPERTY:
		// Property resolution
		resolve := data.GetResolver(data.RES_PROPERTY)
		value, exists = resolve(nil, attrName)
		if !exists {
			err := fmt.Errorf("Failed to resolve Property: '%s'. Ensure that property is configured in the application.", attrName, )
			logger.Error(err.Error())
			return nil, err
		}
	case data.RES_ENV:
		// Environment resolution
		resolve := data.GetResolver(data.RES_ENV)
		value, exists = resolve(nil, attrName)
		if !exists {
			err := fmt.Errorf("Failed to resolve Environment Variable: '%s'. Ensure that variable is configured.", attrName)
			logger.Error(err.Error())
			return "", err
		}
	default:
		//data.RES_ACTIVITY
		//data.RES_TRIGGER
		//data.RES_SCOPE
		resolve := data.GetResolver(resType)
		if resolve != nil {
			value, exists = resolve(scope, attrName)
		}
		if !exists {
			err := fmt.Errorf("Could not resolve '%s' in the current scope", e.rep)
			logger.Error(err.Error())
			return nil, err
		}
	}

	if path != "" {
		value, err = data.PathGetValue(value, path)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
	}

	return value, nil
}

func NewAssignExpr(assignTo string, value interface{}) data.Expr {

	attrName, attrPath, _ := data.PathDeconstruct(assignTo)
	return &assignExpr{assignAttrName: attrName, assignAttrPath: attrPath, value:value}
}

type assignExpr struct {
	assignAttrName string
	assignAttrPath string
	value    interface{}
}

func (e *assignExpr) Eval(scope data.Scope) (interface{}, error) {

	var err error

	if e.assignAttrPath == "" {
		//simple assignment
		err = scope.SetAttrValue(e.assignAttrName, e.value)
		return nil, err
	}

	attr, exists := scope.GetAttr(e.assignAttrName)

	if !exists {
		return nil, fmt.Errorf("Attribute '%s' does not exists\n", e.assignAttrName)
	}

	//temporary hack
	if attr.Value == nil {
		switch attr.Type {
		case data.OBJECT:
			attr.Value = make(map[string]interface{})
		case data.PARAMS:
			attr.Value = make(map[string]string)
		}
	}

	err = data.PathSetValue(attr.Value, e.assignAttrPath, e.value)
	return nil, err
}