package json

import (
	"encoding/json"
	"sync"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/json/field"

	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestSetArrayObject(t *testing.T) {
	mappingField := field.NewMappingField(false, true, []string{"City[0]", "Array[1]", "id"})
	v, err := SetFieldValueFromString("setvalue", jsonData, mappingField)
	assert.Nil(t, err)
	assert.NotNil(t, v)

	d, err := getValue(v, []string{"City[0]", "Array[1]", "id"}, false, true)
	assert.Nil(t, err)
	assert.Equal(t, "setvalue", d)

}

func printvalue(in interface{}) {
	v, _ := json.Marshal(in)
	fmt.Println(string(v))
}

func TestSetRootChildArray(t *testing.T) {
	mappingField := field.NewMappingField(false, true, []string{"Emails[0]"})
	v, err := SetFieldValueFromString("test-cases@gmail.com", jsonData, mappingField)
	assert.Nil(t, err)
	assert.NotNil(t, v)

	d, err := getValue(v, []string{"Emails[0]"}, false, true)
	assert.Nil(t, err)
	assert.Equal(t, "test-cases@gmail.com", d)
}

func TestSetRootArray(t *testing.T) {
	mappingField := field.NewMappingField(false, true, []string{"[0]", "ss"})
	v, err := SetFieldValue("test-cases@gmail.com", "{}", mappingField)
	assert.Nil(t, err)
	assert.NotNil(t, v)

	d, err := getValue(v, []string{"[0]", "ss"}, false, true)
	assert.Nil(t, err)
	assert.Equal(t, "test-cases@gmail.com", d)

}

func TestSetObject(t *testing.T) {
	mappingField := field.NewMappingField(false, false, []string{"ZipCode"})
	v, err := SetFieldValueFromString("77479", jsonData, mappingField)
	assert.Nil(t, err)
	assert.NotNil(t, v)

	d, err := getValue(v, []string{"ZipCode"}, false, true)
	assert.Nil(t, err)
	assert.Equal(t, "77479", d)
}

func TestSetEmptyField(t *testing.T) {
	mappingField := field.NewMappingField(false, false, []string{"ZipCode"})
	jsond := "{}"
	v, err := SetFieldValueFromString("77479", jsond, mappingField)
	assert.Nil(t, err)
	assert.NotNil(t, v)
	printvalue(v)
	d, err := getValue(v, []string{"ZipCode"}, false, true)
	assert.Nil(t, err)
	assert.Equal(t, "77479", d)
}

func TestSetEmptyField4(t *testing.T) {
	jsond := "{}"
	mappingField := field.NewMappingField(false, true, []string{"ZipCode[1]"})
	v, err := SetFieldValueFromString("77479", jsond, mappingField)
	assert.Nil(t, err)
	assert.NotNil(t, v)
	d, err := getValue(v, []string{"ZipCode[1]"}, false, true)
	assert.Nil(t, err)
	assert.Equal(t, "77479", d)
}

func TestSetEmptyField5(t *testing.T) {
	jsond := "{}"
	mappingField := field.NewMappingField(false, true, []string{"ZipCode[1]"})
	v, err := SetFieldValueFromString("77479", jsond, mappingField)

	mappingField = field.NewMappingField(false, true, []string{"ZipCode[0]"})
	v2, err := SetFieldValue("77479", v, mappingField)
	assert.Nil(t, err)
	assert.NotNil(t, v2)
	d, err := getValue(v, []string{"ZipCode[0]"}, false, true)
	assert.Nil(t, err)
	assert.Equal(t, "77479", d)

	d, err = getValue(v, []string{"ZipCode[1]"}, false, true)
	assert.Nil(t, err)
	assert.Equal(t, "77479", d)
}

func TestSetEmptyNestField1(t *testing.T) {
	jsond := "{}"
	mappingField := field.NewMappingField(false, true, []string{"pet", "photoUrls[0]"})
	v, err := SetFieldValueFromString("url", jsond, mappingField)

	mappingField = field.NewMappingField(false, true, []string{"pet", "photoUrls[1]"})

	v, err = SetFieldValue("url2", v, mappingField)
	assert.Nil(t, err)
	assert.NotNil(t, v)

	d, err := getValue(v, []string{"pet", "photoUrls[0]"}, false, true)
	assert.Nil(t, err)
	assert.Equal(t, "url", d)
	d, err = getValue(v, []string{"pet", "photoUrls[1]"}, false, true)
	assert.Nil(t, err)
	assert.Equal(t, "url2", d)
}

func TestNameWithSpace(t *testing.T) {
	jsond := "{}"
	mappingField := field.NewMappingField(false, true, []string{"pet name", "photoUrls[0]"})

	v, err := SetFieldValue("url", jsond, mappingField)
	mappingField = field.NewMappingField(false, true, []string{"pet name", "photoUrls[1]"})

	v, err = SetFieldValue("url2", v, mappingField)
	assert.Nil(t, err)
	assert.NotNil(t, v)

	d, err := getValue(v, []string{"pet name", "photoUrls[0]"}, false, true)
	assert.Nil(t, err)
	assert.Equal(t, "url", d)
	d, err = getValue(v, []string{"pet name", "photoUrls[1]"}, false, true)
	assert.Nil(t, err)
	assert.Equal(t, "url2", d)

}

func TestNameNest2(t *testing.T) {
	jsond := "{}"
	mappingField := field.NewMappingField(false, true, []string{"input", "Account", "records[0]", "ID"})
	v, err := SetFieldValue("id22", jsond, mappingField)

	mappingField = field.NewMappingField(false, true, []string{"input", "Account", "records[0]", "Name"})

	v, err = SetFieldValue("namesssss", v, mappingField)
	assert.Nil(t, err)
	assert.NotNil(t, v)

	d, err := getValue(v, []string{"input", "Account", "records[0]", "ID"}, false, true)
	assert.Nil(t, err)
	assert.Equal(t, "id22", d)
	d, err = getValue(v, []string{"input", "Account", "records[0]", "Name"}, false, true)
	assert.Nil(t, err)
	assert.Equal(t, "namesssss", d)
}

func TestNameSameLevel(t *testing.T) {
	jsond := "{}"
	mappingField := field.NewMappingField(false, true, []string{"input", "Account", "ID"})
	v, err := SetFieldValue("id", jsond, mappingField)

	mappingField = field.NewMappingField(false, true, []string{"input", "Account", "Name"})

	v, err = SetFieldValue("namesssss", v, mappingField)
	assert.Nil(t, err)
	assert.NotNil(t, v)

	d, err := getValue(v, []string{"input", "Account", "ID"}, false, true)
	assert.Nil(t, err)
	assert.Equal(t, "id", d)
	d, err = getValue(v, []string{"input", "Account", "Name"}, false, true)
	assert.Nil(t, err)
	assert.Equal(t, "namesssss", d)

}

func TestNameWithTag(t *testing.T) {
	jsond := "{}"
	mappingField := field.NewMappingField(false, true, []string{"pet", "pet name", "photo	Urls[0]"})
	v, err := SetFieldValue("url", jsond, mappingField)

	mappingField = field.NewMappingField(false, true, []string{"pet", "pet name", "photo	Urls[1]"})

	v, err = SetFieldValue("url2", v, mappingField)
	assert.Nil(t, err)
	assert.NotNil(t, v)

	d, err := getValue(v, []string{"pet", "pet name", "photo	Urls[0]"}, false, true)
	assert.Nil(t, err)
	assert.Equal(t, "url", d)
	d, err = getValue(v, []string{"pet", "pet name", "photo	Urls[1]"}, false, true)
	assert.Nil(t, err)
	assert.Equal(t, "url2", d)
}

func TestSetEmptyNestField(t *testing.T) {
	jsond := "{}"
	mappingField := field.NewMappingField(false, true, []string{"Response", "Pet", "Tags[0]", "Name"})
	v, err := SetFieldValueFromString("tagID", jsond, mappingField)

	mappingField = field.NewMappingField(false, true, []string{"Response", "Pet", "Tags[1]", "Name"})

	v, err = SetFieldValue("tagID2", v, mappingField)

	assert.Nil(t, err)
	assert.NotNil(t, v)
	d, err := getValue(v, []string{"Response", "Pet", "Tags[0]", "Name"}, false, true)
	assert.Nil(t, err)
	assert.Equal(t, "tagID", d)
	d, err = getValue(v, []string{"Response", "Pet", "Tags[1]", "Name"}, false, true)
	assert.Nil(t, err)
	assert.Equal(t, "tagID2", d)
}

func TestConcurrentSet(t *testing.T) {
	w := sync.WaitGroup{}
	var recovered interface{}
	//Create factory

	for r := 0; r < 100000; r++ {
		w.Add(1)
		go func(i int) {
			defer w.Done()
			defer func() {
				if r := recover(); r != nil {
					recovered = r
				}
			}()
			jsond := "{}"
			mappingField := field.NewMappingField(false, true, []string{"pet name", "photoUrls[0]"})
			v, err := SetFieldValue("url", jsond, mappingField)

			mappingField = field.NewMappingField(false, true, []string{"pet name", "photoUrls[1]"})

			v, err = SetFieldValue("url2", v, mappingField)
			assert.Nil(t, err)
			assert.NotNil(t, v)
		}(r)

	}
	w.Wait()
	assert.Nil(t, recovered)
}

func getValue(value interface{}, fields []string, hasSpec, hasArray bool) (interface{}, error) {
	mapField := field.NewMappingField(hasSpec, hasArray, fields)
	return GetFieldValueFromIn(value, mapField)
}
