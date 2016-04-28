package fglua

import (
	"testing"
	"github.com/TIBCOSoftware/flogo-lib/core/flow"
	"encoding/json"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"fmt"
)

const defJSON = `
{
  "type": 1,
  "name": "Demo Flow",
  "model": "simple",
  "rootTask": {
    "id": 1,
    "type": 1,
    "activityType": "",
    "name": "root",
    "tasks": [
      {
        "id": 2,
        "type": 1,
        "name": "A",
        "activityType": ""
      },
      {
        "id": 3,
        "type": 1,
        "name": "B",
        "activityType": ""
      },
      {
        "id": 4,
        "type": 1,
        "name": "C",
        "activityType": ""
      }
    ],
    "links": [
      { "id": 1, "type": 1,  "name": "", "to": 3,  "from": 2, "value":"$[A1.sensor].petId > 50" },
      { "id": 2, "type": 1, "name": "", "to": 4, "from": 2, "value":"$petId > 5" }
    ]
  }
}
`

func TestLuaLinkExprManager_EvalLinkExpr(t *testing.T) {

	defRep := &flow.DefinitionRep{}
	json.Unmarshal([]byte(defJSON), defRep)

	def,_ := flow.NewDefinition(defRep)

	mgr := NewLuaLinkExprManager(def)

	link1 := def.GetLink(1)
	link2 := def.GetLink(2)

	sensorData := make(map[string]interface{})
	sensorData["temp"] = 55

	attrs := []*data.Attribute{
		&data.Attribute{Name:"petId", Type:"integer", Value:3},
		&data.Attribute{Name:"sensorData", Type:"object", Value:sensorData},
	}

	scope := data.NewSimpleScope(attrs, nil)

	result := mgr.EvalLinkExpr(link1, scope)
	fmt.Printf("Link 1 Result: %v\n", result)

	result = mgr.EvalLinkExpr(link2, scope)
	fmt.Printf("Link 2 Result: %v\n", result)

	scope.SetAttrValue("petId",6)
	result = mgr.EvalLinkExpr(link2, scope)

	fmt.Printf("Link2 Result: %v\n", result)
}