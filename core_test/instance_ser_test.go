package core_test

import (
	"encoding/json"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/process"
	"github.com/TIBCOSoftware/flogo-lib/core/processinst"
	"github.com/TIBCOSoftware/flogo-lib/util"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("core")

const defJSON = `
{
    "type": 1,
    "name": "test",
    "model": "test",
    "rootTask": {
      "id": 1,
      "type": 1,
      "activityType": "",
      "name": "root",
      "tasks": [
        {
          "id": 2,
          "type": 1,
          "name": "a"
        },
        {
          "id": 3,
          "type": 1,
          "name": "b"
        }
      ],
      "links": [
        { "id": 1, "type": 1,  "name": "","from": 2, "to": 3 }
      ]
    }
  }
`

func TestFullSerialization(t *testing.T) {

	defRep := &process.DefinitionRep{}
	json.Unmarshal([]byte(defJSON), defRep)

	log.Infof("Def Rep: %v", defRep)

	def := process.NewDefinition(defRep)

	idGen, _ := util.NewGenerator()
	id := idGen.NextAsString()

	instance := processinst.NewProcessInstance(id, "uri1", def)

	instance.Start(nil)

	hasWork := true

	for hasWork && instance.Status() < processinst.StatusCompleted {
		hasWork = instance.DoStep()

		json, _ := json.Marshal(instance)
		log.Debugf("Snapshot: %s\n", string(json))
	}
}

func TestIncrementalSerialization(t *testing.T) {

	defRep := &process.DefinitionRep{}
	json.Unmarshal([]byte(defJSON), defRep)

	idGen, _ := util.NewGenerator()
	id := idGen.NextAsString()

	def := process.NewDefinition(defRep)

	instance := processinst.NewProcessInstance(id, "uri2", def)

	instance.Start(nil)

	hasWork := true

	for hasWork && instance.Status() < processinst.StatusCompleted {
		hasWork = instance.DoStep()

		json, _ := json.Marshal(instance.GetChanges())
		log.Debugf("Changes: %s\n", string(json))
	}
}
