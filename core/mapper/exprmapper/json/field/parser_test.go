package field

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParser(t *testing.T) {
	jsonPath := `Object.Maps3["dd*cc"]["y.x"]["d.d"][110]`
	f, err := ParseMappingField(jsonPath)
	assert.Nil(t, err)
	v, _ := json.Marshal(f.Getfields())
	fmt.Println(string(v))
}
