package json

import (
	"strconv"
	"strings"
	"sync"
)

func getArrayFieldName(fields []string) ([]string, int, int) {
	var tmpFields []string
	index := -1
	var arrayIndex int
	for i, field := range fields {
		if strings.Index(field, "[") >= 0 && strings.Index(field, "]") >= 0 {
			arrayIndex, _ = getFieldSliceIndex(field)
			fieldName := getFieldName(field)
			index = i
			if fieldName != "" {
				tmpFields = append(tmpFields, getFieldName(field))
			}
			break
		} else {
			tmpFields = append(tmpFields, field)
		}
	}
	return tmpFields, index, arrayIndex
}

func hasArrayFieldInArray(fields []string) bool {
	for _, field := range fields {
		if strings.Index(field, "[") >= 0 && strings.HasSuffix(field, "]") {
			//Make sure the index are integer
			_, err := strconv.Atoi(getNameInsideBrancket(field))
			if err == nil {
				return true
			}
		}
	}
	return false
}

func handleSetValue(value interface{}, jsonData *JSONData, fields []string) error {

	log.Debugf("All fields %+v", fields)
	jsonData.rw.Lock()
	defer jsonData.rw.Unlock()

	container := jsonData.container
	if hasArrayFieldInArray(fields) {
		arrayFields, fieldNameindex, arrayIndex := getArrayFieldName(fields)
		//No array field found
		if fieldNameindex == -1 {
			if arrayIndex == -2 {
				//Append
				err := container.ArrayAppend(value, arrayFields...)
				if err != nil {
					return err
				}
			} else {
				//set to exist index array
				size, err := container.ArrayCount(arrayFields...)
				if err != nil {
					return err
				}
				if arrayIndex > size-1 {
					err := container.ArrayAppend(value, arrayFields...)
					if err != nil {
						return err
					}
				} else {
					array := container.S(arrayFields...)
					_, err := array.SetIndex(value, arrayIndex)
					if err != nil {
						return err
					}
				}
			}
		} else {
			restFields := fields[fieldNameindex+1:]
			if arrayFields != nil {
				if container.Exists(arrayFields...) {
					if restFields == nil || len(restFields) <= 0 {
						array, ok := container.Search(arrayFields...).Data().([]interface{})
						if ok {
							if arrayIndex > len(array)-1 {
								array = append(array, value)
							} else {
								array[arrayIndex] = value
							}
						}
						_, err := container.Set(array, arrayFields...)
						return err
					} else {
						var element *Container
						var err error
						count, err := container.ArrayCount(arrayFields...)
						if err != nil {
							return err
						}
						if arrayIndex > count-1 {
							maps := make(map[string]interface{})
							newObject, _ := Consume(maps)
							_, err = newObject.Set(value, restFields...)
							log.Debugf("new object %s", newObject.String())
							if err != nil {
								return err
							}

							err = container.ArrayAppend(maps, arrayFields...)
							if err != nil {
								return err
							}
							if !hasArrayFieldInArray(restFields) {
								return nil
							}

						}

						element, err = container.ArrayElement(arrayIndex, arrayFields...)
						if err != nil {
							return err
						}
						return handleSetValue(value, &JSONData{container: element, rw: sync.RWMutex{}}, restFields)
					}
				}

			} else if fieldNameindex == 0 && getFieldName(fields[fieldNameindex]) == "" {
				//Root only [0]
				array, ok := container.Data().([]interface{})
				if !ok {

					toArray, err := ToArray(container.Data())
					if err != nil {
						toArray = make([]interface{}, arrayIndex+1)
					}

					array = toArray
				}
				container.object = array
				if restFields == nil || len(restFields) <= 0 {
					_, err := container.SetIndex(value, arrayIndex)
					return err
				}
			}
			//Create new one
			array, err := container.ArrayOfSize(arrayIndex+1, arrayFields...)
			if err != nil {
				return err
			}

			if hasArrayFieldInArray(restFields) {
				return handleSetValue(value, &JSONData{container: array, rw: sync.RWMutex{}}, restFields)
			}
			maps := make(map[string]interface{})
			newObject, _ := Consume(maps)
			_, err = newObject.Set(value, restFields...)
			log.Debugf("new object %s", newObject.String())
			if err != nil {
				return err
			}
			_, err = array.SetIndex(newObject.object, arrayIndex)
		}
	} else {
		_, err := jsonData.container.Set(value, fields...)
		if err != nil {
			return err
		}
	}
	return nil
}

func handleGetValue(jsonData *JSONData, fields []string) (interface{}, error) {

	log.Debugf("All fields %+v", fields)
	jsonData.rw.Lock()
	defer jsonData.rw.Unlock()

	container := jsonData.container
	if hasArrayFieldInArray(fields) {
		arrayFields, fieldNameindex, arrayIndex := getArrayFieldName(fields)
		//No array field found
		if fieldNameindex == -1 {
			return container.S(arrayFields...).Data(), nil
		}
		restFields := fields[fieldNameindex+1:]
		specialField, err := container.ArrayElement(arrayIndex, arrayFields...)
		if err != nil {
			return nil, err
		}
		log.Debugf("Array element value %s", specialField)
		if hasArrayFieldInArray(restFields) {
			return handleGetValue(&JSONData{container: specialField, rw: sync.RWMutex{}}, restFields)
		}
		return specialField.S(restFields...).Data(), nil
	}
	log.Debugf("No array found for array %+v and size %d", fields, len(fields))
	return container.S(fields...).Data(), nil
}
