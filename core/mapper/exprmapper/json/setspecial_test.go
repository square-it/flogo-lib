package json

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/json/field"

	"github.com/stretchr/testify/assert"
)

func TestSetSpecialObjectField(t *testing.T) {
	// path := `Object.Maps3["dd*cc"]["y.x"][d.d]`
	mappingField := field.NewMappingField(true, false, []string{"Object", "Maps3", "dd*cc", "y.x", "d.d"})

	value, err := SetFieldValueFromString("lixi", "{}", mappingField)
	assert.Nil(t, err)
	assert.NotNil(t, value)
	log.Info(value)
	v, _ := json.Marshal(value)
	assert.Equal(t, `{"Object":{"Maps3":{"dd*cc":{"y.x":{"d.d":"lixi"}}}}}`, string(v))
}

func TestSetSpecialArrayField2(t *testing.T) {
	// path := `Object.Maps3[0]["dd*cc"]`
	mappingField := field.NewMappingField(true, true, []string{"Object", "Maps3[0]", "dd*cc"})

	value, err := SetFieldValueFromString("lixi", "{}", mappingField)
	assert.Nil(t, err)
	assert.NotNil(t, value)
	log.Info(value)
	v, _ := json.Marshal(value)
	fmt.Println(string(v))
	assert.Equal(t, `{"Object":{"Maps3":[{"dd*cc":"lixi"}]}}`, string(v))
}

func TestSetSpecialArrayFieldMultipleLEvel(t *testing.T) {
	// path := `Object.Maps3["dd.cc"][0]["y.x"][d.d].name`
	mappingField := field.NewMappingField(true, true, []string{"Object", "Maps3", "dd.cc[0]", "y.x", "d.d", "name"})
	value, err := SetFieldValueFromString("lixi", "{}", mappingField)
	assert.Nil(t, err)
	assert.NotNil(t, value)
	log.Info(value)
	v, _ := json.Marshal(value)
	fmt.Println(string(v))
	assert.Equal(t, `{"Object":{"Maps3":{"dd.cc":[{"y.x":{"d.d":{"name":"lixi"}}}]}}}`, string(v))
}

func TestSetArrayRootOnly(t *testing.T) {
	mappingField := field.NewMappingField(true, true, []string{"[0]"})
	value, err := SetFieldValueFromString("lixi", "{}", mappingField)
	assert.Nil(t, err)
	assert.NotNil(t, value)
	v, err := getValue(value, []string{"[0]"}, false, true)
	assert.Nil(t, err)
	assert.Equal(t, "lixi", v)

	mappingField = field.NewMappingField(true, true, []string{"[0]", "tt"})
	value, err = SetFieldValueFromString("ssssss", "{}", mappingField)
	assert.Nil(t, err)
	assert.NotNil(t, value)
	v, err = getValue(value, []string{"[0]", "tt"}, false, true)
	assert.Nil(t, err)
	assert.Equal(t, "ssssss", v)
}
