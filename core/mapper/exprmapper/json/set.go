package json

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/json/field"
)

type JSONData struct {
	container *Container
	rw        sync.RWMutex
}

//func SetFieldValueFromStringP(data interface{}, jsonData string, path string) (interface{}, error) {
//	jsonParsed, err := ParseJSON([]byte(jsonData))
//	if err != nil {
//		return nil, err
//
//	}
//	return setValueP(data, &JSONData{container: jsonParsed, rw: sync.RWMutex{}}, path)
//}

func SetFieldValueFromString(data interface{}, jsonData string, mappingField *field.MappingField) (interface{}, error) {
	jsonParsed, err := ParseJSON([]byte(jsonData))
	if err != nil {
		return nil, err

	}
	container := &JSONData{container: jsonParsed, rw: sync.RWMutex{}}
	err = setValue(data, container, mappingField)
	return container.container.object, err
}

//
//func SetFieldValueP(data interface{}, jsonData interface{}, path string) (interface{}, error) {
//	switch t := jsonData.(type) {
//	case string:
//		return SetFieldValueFromStringP(data, t, path)
//	default:
//		jsonParsed, err := Consume(jsonData)
//		if err != nil {
//			return nil, err
//
//		}
//		return setValueP(data, &JSONData{container: jsonParsed, rw: sync.RWMutex{}}, path)
//	}
//}

func SetFieldValue(data interface{}, jsonData interface{}, mappingField *field.MappingField) (interface{}, error) {
	switch t := jsonData.(type) {
	case string:
		return SetFieldValueFromString(data, t, mappingField)
	default:
		jsonParsed, err := Consume(jsonData)
		if err != nil {
			return nil, err

		}
		container := &JSONData{container: jsonParsed, rw: sync.RWMutex{}}
		setValue(data, container, mappingField)
		return container.container.object, nil
	}
}

//func setValueP(value interface{}, jsonData *JSONData, mapField *field.MappingField) (interface{}, error) {
//	if mapField.HasArray() && mapField.HasSepcialField() {
//		return handleArrayWithSpecialFields(value, jsonData, mapField.Getfields())
//	} else if mapField.HasArray() {
//		return setArrayValue(value, jsonData, strings.Join(mapField.Getfields(), "."))
//	} else if mapField.HasSepcialField() {
//		_, err := jsonData.container.Set(value, mapField.Getfields()...)
//		if err != nil {
//			return nil, err
//		}
//		return jsonData.container.object, nil
//	}
//	_, err := jsonData.container.Set(value, mapField.Getfields()...)
//	if err != nil {
//		return nil, err
//	}
//	return jsonData.container.object, nil
//}

func setValue(value interface{}, jsonData *JSONData, mappingField *field.MappingField) error {
	if mappingField.HasArray() && mappingField.HasSepcialField() {
		return handleArrayWithSpecialFields(value, jsonData, mappingField.Getfields())
	} else if mappingField.HasArray() {
		return handleArrayWithSpecialFields(value, jsonData, mappingField.Getfields())
	} else if mappingField.HasSepcialField() {
		_, err := jsonData.container.Set(value, mappingField.Getfields()...)
		if err != nil {
			return err
		}
		return nil
	}
	_, err := jsonData.container.Set(value, mappingField.Getfields()...)
	return err
}

func setArrayValue(value interface{}, jsonData *JSONData, path string) error {

	jsonData.rw.Lock()

	defer jsonData.rw.Unlock()

	container := jsonData.container
	if strings.Index(path, "[") >= 0 {
		//Array, if fieldname empty then take it as root array
		arrayFieldName := getFieldName(path)
		index, _ := getFieldSliceIndex(path)

		if arrayFieldName == "" {
			//No field name to with root field
			//check if the root array already exist.
			array, ok := container.Data().([]interface{})
			if !ok {

				toArray, err := ToArray(container.Data())
				if err != nil {
					toArray = make([]interface{}, index+1)
					//for i, v := range toArray {
					//	toArray[i] =
					//}
				}

				array = toArray
			}
			container.object = array
		}
		log.Debug("Field Name:", arrayFieldName, " and index: ", index)
		restPath := getRestArrayFieldName(path)
		if restPath == "" {
			if strings.Index(path, "]") == len(path)-1 {
				if arrayFieldName == "" {
					container.SetIndex(value, index)
					return nil
				}
				if container.ExistsP(arrayFieldName) {
					if index == -2 {
						//Append
						err := container.ArrayAppend(value, strings.Split(arrayFieldName, ".")...)
						if err != nil {
							return err
						}
					} else {
						//set to exist index array
						size, err := container.ArrayCountP(arrayFieldName)
						if err != nil {
							return err
						}
						if index > size-1 {
							err := container.ArrayAppendP(value, arrayFieldName)
							if err != nil {
								return err
							}
						} else {
							array := container.Path(arrayFieldName)
							_, err := array.SetIndex(value, index)
							if err != nil {
								return err
							}
						}
					}

				} else {
					//Not exist so init a new array
					if index == -2 {
						_, err := container.Array(strings.Split(arrayFieldName, ".")...)
						if err != nil {
							return err
						}
						err = container.ArrayAppend(value, strings.Split(arrayFieldName, ".")...)
						if err != nil {
							return err
						}
					} else {
						//Since make array with index lengh
						array, err := container.ArrayOfSize(index+1, strings.Split(arrayFieldName, ".")...)
						if err != nil {
							return err
						}
						_, err = array.SetIndex(value, index)
						if err != nil {
							return err
						}
					}
				}

			} else {
				jsonField := container.Path(arrayFieldName)
				_, err := jsonField.SetIndex(value, index)
				if err != nil {
					return err
				}
			}

		} else {
			if arrayFieldName == "" {
				c := container.Index(index)
				data := &JSONData{container: c, rw: sync.RWMutex{}}
				err := setArrayValue(value, data, restPath)
				if err != nil {
					return err
				}
				_, err = container.SetIndex(data.container.object, index)
				return err
			}
			if container.ExistsP(arrayFieldName) {
				size, err := container.ArrayCountP(arrayFieldName)
				if err != nil {
					return err
				}

				if index > size-1 {

					newObject, err := ParseJSON([]byte("{}"))
					_, err = newObject.SetP(value, restPath)
					log.Debugf("new object %s", newObject.String())
					if err != nil {
						return err
					}
					//o ,_ := ParseJSON(newObject.Bytes())
					maps := &map[string]interface{}{}
					err = json.Unmarshal(newObject.Bytes(), maps)
					if err != nil {
						return err
					}

					err = container.ArrayAppendP(maps, arrayFieldName)
					if err != nil {
						return err
					}

					if strings.Index(restPath, "[") > 0 {
						//TODO
						c, err := container.ArrayElementP(index, arrayFieldName)
						if err != nil {
							return err
						}
						return setArrayValue(value, &JSONData{container: c, rw: sync.RWMutex{}}, restPath)
					} else {
						//_, err := jsonField.Set(value, restPath)
						//if err != nil {
						//	return  err
						//}
					}
				} else {

					jsonField, err := container.ArrayElementP(index, arrayFieldName)
					//arraySize
					if err != nil {
						return err
					}
					if strings.Index(restPath, "[") > 0 {
						return setArrayValue(value, &JSONData{container: jsonField, rw: sync.RWMutex{}}, restPath)
					} else {
						switch t := jsonField.object.(type) {
						case map[string]interface{}:
							jsonField.object = t
						case *map[string]interface{}:
							jsonField.object = *t
						}
						_, err := jsonField.SetP(value, restPath)
						if err != nil {
							return err
						}
					}
				}

			} else {
				//Not exist so init a new array
				//Since make array with index lengh
				array, err := container.ArrayOfSize(index+1, strings.Split(arrayFieldName, ".")...)
				if err != nil {
					return err
				}

				if strings.Index(restPath, "[") > 0 {
					return setArrayValue(value, &JSONData{container: array, rw: sync.RWMutex{}}, restPath)
				} else {
					newObject, err := ParseJSON([]byte("{}"))
					_, err = newObject.SetP(value, restPath)
					log.Debugf("new object %s", newObject.String())
					if err != nil {
						return err
					}
					//o ,_ := ParseJSON(newObject.Bytes())
					maps := &map[string]interface{}{}
					err = json.Unmarshal(newObject.Bytes(), maps)
					if err != nil {
						return err
					}
					_, err = array.SetIndex(maps, index)
				}
			}

		}
		// }
	} else {
		_, err := container.Set(value, strings.Split(path, ".")...)
		if err != nil {
			return err
		}

	}
	return nil
}
