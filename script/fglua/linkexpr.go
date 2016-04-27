package fglua

import (
	"github.com/TIBCOSoftware/flogo-lib/core/flow"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/Shopify/go-lua"
	"encoding/json"
	"fmt"
	"strings"
	"bytes"
	"strconv"
)

type LuaLinkExprManager struct {
	Values map[int][]string
	L      *lua.State
}

func NewLuaLinkExprManager(def *flow.Definition) *LuaLinkExprManager{

	mgr := &LuaLinkExprManager{}
	mgr.Values = make(map[int][]string)

	links := flow.GetExpressionLinks(def)

	fmt.Printf("links: %v", links)

	var buffer bytes.Buffer

	for _, link := range links {
		attrs, expr := transExpr(link.Value())

		mgr.Values[link.ID()] = attrs

		buffer.WriteString("l")
		buffer.WriteString(strconv.Itoa(link.ID()))
		buffer.WriteString(" = function (v) \n return ")
		buffer.WriteString(expr)
		buffer.WriteString("\nend\n")

		fmt.Println(expr)

	}

	script := buffer.String()

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
			attrs = append(attrs, s[i + 1:j])
			rattrs = append(rattrs, s[i:j])
			rattrs = append(rattrs, `v["` + s[i + 1:j] + `"]`)
			i = j
		}
	}

	replacer := strings.NewReplacer(rattrs...)

	return attrs, replacer.Replace(s)
}

func (em *LuaLinkExprManager) EvalLinkExpr(link *flow.Link, scope data.Scope) bool {

	if link.Type() == flow.LtDependency {
		// dependency links are always true
		return true
	}

	attrs := em.Values[link.ID()]

	vals := make(map[string]interface{})

	for _, attr := range attrs {
		val,_ := scope.GetAttrValue(attr)
		vals[attr] = val
	}

	em.L.Global("l"+strconv.Itoa(link.ID()))
	PushMap(em.L, vals)
	em.L.Call(1, 1)
	ret := em.L.ToValue(-1)

	return ret.(bool)
}

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

