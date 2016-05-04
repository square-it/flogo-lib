package fglua

import (
	"sync"
	"github.com/Shopify/go-lua"
	"encoding/json"
)

type LuaStatePool struct {
	mutex sync.Mutex
	saved []*lua.State
}

func (pl *LuaStatePool) Get() (*lua.State, bool) {
	pl.mutex.Lock()
	defer pl.mutex.Unlock()
	n := len(pl.saved)
	if n == 0 {
		return lua.NewState(), true
	}
	x := pl.saved[n-1]
	pl.saved = pl.saved[0 : n-1]
	return x, false
}

func (pl *LuaStatePool) new() *lua.State {
	L := lua.NewState()
	// setting the L up here.
	// load scripts, set global variables, share channels, etc...
	return L
}

func (pl *LuaStatePool) Release(L *lua.State) {
	pl.mutex.Lock()
	defer pl.mutex.Unlock()
	pl.saved = append(pl.saved, L)
}

func (pl *LuaStatePool) Shutdown() {
	pl.saved = nil
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
