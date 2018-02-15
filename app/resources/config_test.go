package resources

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const resJSON = `
{
  "resources":
  [
    {
      "type": "flow",
      "entries":[
      {
        "id":"myflow"
      }
      ]
    },
    {
      "type": "schema",
      "entries": [
      {
        "id":"myschema1"
      }
      ]
    },
    {
      "type": "sharedconfig",
      "entries":[
      {
        "id":"myconfiguration1"
      }
      ]
    }
  ]
}
`

func TestDeserialize(t *testing.T) {

	defRep := &ResourcesConfig{}

	err := json.Unmarshal([]byte(resJSON), defRep)
	assert.Nil(t, err)

	fmt.Printf("Resources: %v", defRep)
}

