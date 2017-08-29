package json

import (
	"encoding/json"
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

func TestSetArrayObject(t *testing.T) {
	v, err := SetFieldValueFromString("4444555", jsonData, "City[0].Array[1].id")
	assert.Nil(t, err)
	assert.NotNil(t, v)
	logger.Debug("FInaly value:", v)

}

func TestSetRootArray(t *testing.T) {
	v, err := SetFieldValueFromString("lixingwang@gmail.com", jsonData, "Emails[0]")
	assert.Nil(t, err)
	assert.NotNil(t, v)
	logger.Debug("FInaly value:", v)

}

func TestSetObject(t *testing.T) {
	v, err := SetFieldValueFromString("77479", jsonData, "ZipCode")
	assert.Nil(t, err)
	assert.NotNil(t, v)
	logger.Debug("FInaly value:", v)
}

func TestSetEmptyField(t *testing.T) {
	jsond := "{}"
	v, err := SetFieldValueFromString("77479", jsond, "ZipCode")
	assert.Nil(t, err)
	assert.NotNil(t, v)
	logger.Debug("FInaly value:", v)
}

func TestSetEmptyField3(t *testing.T) {
	jsond := "{}"
	v, err := SetFieldValueFromString("77479", jsond, "ZipCode[]")
	assert.Nil(t, err)
	assert.NotNil(t, v)
	logger.Debug("FInaly value:", v)
}

func TestSetEmptyField4(t *testing.T) {
	jsond := "{}"
	v, err := SetFieldValueFromString("77479", jsond, "ZipCode[1]")
	assert.Nil(t, err)
	assert.NotNil(t, v)
	logger.Debug("FInaly value:", v)
}

func TestSetEmptyField5(t *testing.T) {
	jsond := "{}"
	v, err := SetFieldValueFromString("77479", jsond, "ZipCode[1]")
	v2, err := SetFieldValue("77479", v, "ZipCode[0]")
	assert.Nil(t, err)
	assert.NotNil(t, v2)
	logger.Debug("FInaly value:", v2)
}

func TestSetEmptyArrayField(t *testing.T) {
	jsond := "{}"
	v, err := SetFieldValueFromString("id", jsond, "pet.id")
	logger.Debug("ID value:", v)

	v, err = SetFieldValue("name", v, "pet.name")
	logger.Debug("Name value:", v)

	v, err = SetFieldValue("url", v, "pet.photoUrls[0]")
	assert.Nil(t, err)
	assert.NotNil(t, v)
	logger.Debug("FInaly value:", v)
}

func TestSetEmptyNestField1(t *testing.T) {
	jsond := "{}"
	v, err := SetFieldValueFromString("url", jsond, "pet.photoUrls[0]")
	logger.Debug("First T ", v)
	v, err = SetFieldValue("url2", v, "pet.photoUrls[1]")
	assert.Nil(t, err)
	assert.NotNil(t, v)
	logger.Debug("FInaly value:", v)
}

func TestNameWithSpace(t *testing.T) {
	jsond := "{}"
	v, err := SetFieldValue("url", jsond, "pet name.photoUrls[0]")
	logger.Debug("First T ", v)
	v, err = SetFieldValue("url2", v, "pet name.photoUrls[1]")
	assert.Nil(t, err)
	assert.NotNil(t, v)
	vv, _ := json.Marshal(v)
	logger.Info("FInaly value:", string(vv))
}

func TestNameNest2(t *testing.T) {
	jsond := "{}"
	v, err := SetFieldValue("id", jsond, "input.Account.records[0].ID")
	logger.Debug("First T ", v)
	v, err = SetFieldValue("namesssss", v, "input.Account.records[0].Name")
	assert.Nil(t, err)
	assert.NotNil(t, v)
	vv, _ := json.Marshal(v)
	logger.Info("FInaly value:", string(vv))
}

func TestNameSameLevel(t *testing.T) {
	jsond := "{}"
	v, err := SetFieldValue("id", jsond, "input.Account.ID")
	logger.Debug("First T ", v)
	v, err = SetFieldValue("namesssss", v, "input.Account.Name")
	assert.Nil(t, err)
	assert.NotNil(t, v)
	vv, _ := json.Marshal(v)
	logger.Info("FInaly value:", string(vv))
}

func TestNameWithTag(t *testing.T) {
	jsond := "{}"
	v, err := SetFieldValue("url", jsond, "pet name.photo	Urls[0]")
	logger.Debug("First T ", v)
	v, err = SetFieldValue("url2", v, "pet name.photo	Urls[1]")
	assert.Nil(t, err)
	assert.NotNil(t, v)
	logger.Info("FInaly value:", v)
}

func TestSetEmptyNestField(t *testing.T) {
	jsond := "{}"
	v, err := SetFieldValueFromString("tagID", jsond, "Response.Pet.Tags[0].Name")
	logger.Debug("First T ", v)
	v, err = SetFieldValue("tagID2", v, "Response.Pet.Tags[1].Name")

	assert.Nil(t, err)
	assert.NotNil(t, v)
	logger.Debug("FInaly value:", v)
}

func TestMap(t *testing.T) {
	maps := map[string]interface{}{}

	mapTags := map[string]interface{}{}

	nameMap := map[string]interface{}{}
	nameMap["Name"] = "tagID"
	mapTags["Tags"] = []interface{}{nameMap}

	petTags := map[string]interface{}{}
	petTags["Pet"] = mapTags
	maps["Response"] = petTags
	//map[Response:map[Pet:map[Tags:[{"Name":"tagID"}]]]]
	//map[Response:map[Pet:map[Tags:[map[Name:tagID]]]]]
	//map[Response:map[Pet:map[Tags:[map[]]]]]
	logger.Debug(maps)

	v, _ := json.Marshal(maps)
	logger.Debug(string(v))

}

func setArray(value interface{}, index int, path string) {
	jsonParsed, err := ParseJSON([]byte(jsonData))
	array := jsonParsed.Path(path)

	c, err := array.SetIndex(value, index)
	if err != nil {
		panic(err)
	}

	logger.Debug("Set Value :", c)
	logger.Debug("Final Data:", jsonParsed.String())

}

func getArray(index int, path string) {
	jsonParsed, err := ParseJSON([]byte(jsonData))
	//array := jsonParsed.Path(path)

	c, err := jsonParsed.ArrayElement(index, path)
	if err != nil {
		panic(err)
	}
	logger.Debug("Type :", reflect.TypeOf(c))

	logger.Debug("Get Value :", c.String())

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
			v, err := SetFieldValue("url", jsond, "pet name.photoUrls[0]")
			logger.Debug("First T ", v)
			v, err = SetFieldValue("url2", v, "pet name.photoUrls[1]")
			assert.Nil(t, err)
			assert.NotNil(t, v)
		}(r)

	}
	w.Wait()
	assert.Nil(t, recovered)
}
