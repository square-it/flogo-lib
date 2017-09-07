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
		value, exists = resolve(scope, attrName)
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

func NewAssignExpr(exprRep string, value interface{}) data.Expr {
	return &assignExpr{rep: exprRep, value:value}
}

type assignExpr struct {
	rep string
	value interface{}
}

func (e *assignExpr) Eval(scope data.Scope) (interface{}, error) {

	resType, attrName, path, err := data.GetResolutionInfo(e.rep)

	if err != nil {
		return nil, err
	}

	if resType == data.RES_DEFAULT {
		fmt.Printf("AttrName: %s, Path: %s\n", attrName, path)
	} else if resType != data.RES_SCOPE {
		return nil, fmt.Errorf("Cannot assign to: %s\n", e.rep)
	}

	if path == "" {
		//simple assignment
		err = scope.SetAttrValue(attrName, e.value)
		return nil, err
	}

	attr, exists := scope.GetAttr(attrName)

	if !exists {
		return nil, fmt.Errorf("Attribute '%s' does not exists\n", attrName)
	}

	err = data.PathSetValue(attr.Value, path, e.value)
	return nil, err
}