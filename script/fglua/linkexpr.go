package fglua

import (
	"bytes"
	"encoding/json"
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
	Values map[int][]string
	L      *lua.State
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

	script := buffer.String()
	log.Debugf("Definition [%s] Lua Expressions Script:\n %s\n", def.Name(), script)

	fmt.Println(script)

	mgr.L = lua.NewState()
	lua.DoString(mgr.L, script)

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

	em.L.Global("l" + strconv.Itoa(link.ID()))
	PushMap(em.L, vals)
	em.L.Call(1, 1)
	ret := em.L.ToValue(-1)

	return ret.(bool)
}

// PushVal pushes a value onto the Lua vm's stack
func PushVal(L *lua.State, val interface{}) {
	switch t := val.(type) {
	case string:
		L.PushString(t)
	case int:
		L.PushInteger(t)
	case float64:
		L.PushNumber(t)
	case json.Number:
		f, _ := t.Float64()
		L.PushNumber(f)
	case bool:
		L.PushBoolean(t)
	case nil:
		L.PushNil()
	case map[string]interface{}:
		PushMap(L, t)
	default:
		L.PushUserData(t)
	}
}

// PushMap pushes a map onto the Lua vm's stack
func PushMap(L *lua.State, mapVal map[string]interface{}) int {

	if len(mapVal) > 0 {
		L.CreateTable(0, len(mapVal))

		for k, v := range mapVal {

			PushVal(L, k)
			PushVal(L, v)
			L.SetTable(-3)
		}

	} else {
		L.PushNil()
	}
	return 1
}
