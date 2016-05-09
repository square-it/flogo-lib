package fggos

import (
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/flow"
	"github.com/op/go-logging"
	"github.com/japm/goScript"
)

var log = logging.MustGetLogger("fggos")

// LuaLinkExprManager is the Lua Implementation of a Link Expression Manager
type GosLinkExprManager struct {
	values map[int][]string
	exprs  map[int]*goScript.Expr
}

// NewGosLinkExprManager creates a new LuaLinkExprManager
func NewGosLinkExprManager(def *flow.Definition) *GosLinkExprManager {

	mgr := &GosLinkExprManager{}
	mgr.values = make(map[int][]string)
	mgr.exprs = make(map[int]*goScript.Expr)

	links := flow.GetExpressionLinks(def)

	for _, link := range links {

		if len(strings.TrimSpace(link.Value())) > 0 {
			attrs, exprStr := transExpr(link.Value())

			mgr.values[link.ID()] = attrs

			log.Debugf("expr: %v\n", exprStr)

			expr := &goScript.Expr{}
			err := expr.Prepare(exprStr)

			if err == nil {
				mgr.exprs[link.ID()] = expr
			} else {
				log.Errorf("Error preparing expression: %s - %v", link.Value(), err)
			}
		}
	}

	return mgr
}

func transExpr(s string) ([]string, string) {

	var attrs []string
	var rattrs []string

	strLen := len(s)

	for i := 0; i < strLen; i++ {
		if s[i] == '$' {
			var j int
			for j = i + 1; j < strLen; j++ {
				if s[j] == ' ' {
					break
				}
			}
			attrs = append(attrs, s[i+1:j])
			rattrs = append(rattrs, s[i:j])
			rattrs = append(rattrs, `v["`+s[i+1:j]+`"]`)
			i = j
		}
	}

	replacer := strings.NewReplacer(rattrs...)

	return attrs, replacer.Replace(s)
}

// EvalLinkExpr implements LinkExprManager.EvalLinkExpr
func (em *GosLinkExprManager) EvalLinkExpr(link *flow.Link, scope data.Scope) bool {

	if link.Type() == flow.LtDependency {
		// dependency links are always true
		return true
	}

	attrs, attrsOK := em.values[link.ID()]
	expr, exprOK := em.exprs[link.ID()]

	if !attrsOK || !exprOK {

		log.Warning("Unable to evaluate expression '%s', did not compile properly\n", link.Value())
		return false
	}

	vals := make(map[string]interface{})

	for _, attr := range attrs {

		var attrValue interface{}
		var exists bool

		attrName, attrPath := data.GetAttrPath(attr)

		attrValue, exists = scope.GetAttrValue(attrName)

		if exists && len(attrPath) > 0 {
			//for now assume if we have a path, attr is "object" and only one level
			valMap := attrValue.(map[string]interface{})
			//todo what if the value does not exists
			attrValue, _ = valMap[attrPath]
		}

		vals[attr] = attrValue
	}

	ctxt := make(map[string]interface{})
	ctxt["v"] = vals

	val, _ := expr.Eval(ctxt)
	//todo handle error
	return val.(bool)
}
