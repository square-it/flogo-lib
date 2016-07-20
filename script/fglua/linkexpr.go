package fglua

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/Shopify/go-lua"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/flow"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("fglua")

// LuaLinkExprManager is the Lua Implementation of a Link Expression Manager
type LuaLinkExprManager struct {
	Values  map[int][]string
	script  string
	luaPool *LuaStatePool
}

// NewLuaLinkExprManager creates a new LuaLinkExprManager
func NewLuaLinkExprManager(def *flow.Definition) *LuaLinkExprManager {

	mgr := &LuaLinkExprManager{}
	mgr.Values = make(map[int][]string)

	links := flow.GetExpressionLinks(def)

	var buffer bytes.Buffer

	for _, link := range links {

		if len(strings.TrimSpace(link.Value())) > 0 {
			attrs, expr := transExpr(link.Value())

			mgr.Values[link.ID()] = attrs

			buffer.WriteString("l")
			buffer.WriteString(strconv.Itoa(link.ID()))
			buffer.WriteString(" = function (v) \n return ")
			buffer.WriteString(expr)
			buffer.WriteString("\nend\n")

			log.Debugf("Link[%d] Lua Expression: %s", link.ID(), expr)
			fmt.Println(expr)
		}
	}

	mgr.script = buffer.String()
	log.Debugf("Definition [%s] Lua Expressions Script:\n %s\n", def.Name(), mgr.script)

	mgr.luaPool = &LuaStatePool{saved: make([]*lua.State, 0, 5),}

	return mgr
}

func transExpr(s string) ([]string, string) {

	var attrs []string
	var rattrs []string

	strLen := len(s)

	for i := 0; i < strLen; i++ {
		if s[i] == '$' {

			ignoreBrackets := s[i+1] == '['
			var partOfName bool

			var j int
			for j = i + 1; j < strLen; j++ {

				partOfName, ignoreBrackets = isPartOfName(s[j], ignoreBrackets)

				if !partOfName {
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

func isPartOfName(char byte, ignoreBrackets bool) (bool, bool) {

	if (char < '0' || char > '9') && (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') && char != '.' && char != '_' {

		if  ignoreBrackets && char == '[' {
			return true, true
		} else if ignoreBrackets && char ==']' {
			return true, false
		}

		return false, ignoreBrackets

	}

	return true, ignoreBrackets
}

// EvalLinkExpr implements LinkExprManager.EvalLinkExpr
func (em *LuaLinkExprManager) EvalLinkExpr(link *flow.Link, scope data.Scope) bool {

	if link.Type() == flow.LtDependency {
		// dependency links are always true
		return true
	}

	attrs, ok := em.Values[link.ID()]

	if !ok {
		return false
	}

	vals := make(map[string]interface{})

	for _, attr := range attrs {

		var attrValue interface{}
		var exists bool

		attrName, attrPath, _ := data.GetAttrPath(attr)

		attrValue, exists = scope.GetAttrValue(attrName)

		if exists && len(attrPath) > 0 {
			//for now assume if we have a path, attr is "object" and only one level
			valMap := attrValue.(map[string]interface{})
			//todo what if the value does not exists
			attrValue, _ = valMap[attrPath]
		}

		vals[attr] = attrValue
	}

	L, isNew := em.luaPool.Get()
	defer em.luaPool.Release(L)

	if isNew {
		lua.DoString(L, em.script)
	}

	L.Global("l" + strconv.Itoa(link.ID()))
	PushMap(L, vals)
	L.Call(1, 1)
	ret := L.ToValue(-1)

	return ret.(bool)
}
