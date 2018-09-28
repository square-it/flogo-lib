package ref

import (
	"fmt"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/json/field"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

func (m *MappingRef) GetValueFromOutputScope(mapfield *field.MappingField, outputtscope data.Scope) (interface{}, error) {
	fieldName, err := m.GetFieldName(mapfield)
	if err != nil {
		return nil, err
	}
	log.Debugf("GetValueFromOutputScope field name %s", fieldName)

	attribute, exist := outputtscope.GetAttr(fieldName)
	log.Debugf("GetValueFromOutputScope field name %s and exist %t ", fieldName, exist)

	if exist {
		switch attribute.Type() {
		case data.TypeComplexObject:
			complexObject := attribute.Value().(*data.ComplexObject)
			object := complexObject.Value
			//Convert the object to exist struct.
			//TODO return interface rather than string
			if object == nil {
				return "{}", nil
			}
			return object, nil
		default:
			return attribute.Value(), nil
		}

	}
	return nil, fmt.Errorf("Cannot found attribute %s", fieldName)
}

func (m *MappingRef) GetActivtyRootField(field *field.MappingField) (string, error) {
	if field.HasSepcialField() {
		fields := field.Getfields()

		activityNameRef := fields[0]
		if strings.HasPrefix(activityNameRef, "$") {
			activityName := activityNameRef[1:]
			activityName = getActivityName(activityName)
			fieldName := "_A." + activityName + "." + getFieldName(fields[1])
			return fieldName, nil
		}
		return getFieldName(fields[0]), nil
	}

	if strings.HasPrefix(m.ref, "$") {
		log.Debugf("Mapping ref %s", m.ref)
		mappingFields := strings.Split(m.ref, ".")
		//Which might like $A3.
		//field := mappingFields[1]
		var activityID string
		if strings.HasPrefix(m.ref, "$") {
			activityID = mappingFields[0][1:]
		} else {
			activityID = mappingFields[0]
		}

		//fieldName := "{" + activityID + "." + getFieldName(mappingFields[1]) + "}"
		fieldName := "_A." + getActivityName(activityID) + "." + getFieldName(mappingFields[1])

		log.Debugf("Field name now is: %s", fieldName)
		return fieldName, nil
	} else if strings.Index(m.ref, ".") > 0 {
		log.Debugf("Mapping ref %s", m.ref)
		mappingFields := strings.Split(m.ref, ".")
		log.Debugf("Field name now is: %s", mappingFields[0])
		return getFieldName(mappingFields[0]), nil
	} else {
		return m.ref, nil
	}
}

//
func (m *MappingRef) GetMapToFields(mapField *field.MappingField) (*field.MappingField, error) {
	hasArray := mapField.HasArray()
	fields := mapField.Getfields()
	activityNameRef := fields[0]
	if strings.HasPrefix(activityNameRef, "$") {
		if strings.HasSuffix(fields[1], "]") {
			//Root element is an array
			arrayIndexPart := getArrayIndexPart(fields[1])
			fields[1] = arrayIndexPart
			return field.NewMappingField(mapField.HasSepcialField(), mapField.HasArray(), fields[1:]), nil
		} else {
			return field.NewMappingField(hasArray, mapField.HasArray(), fields[2:]), nil
		}
	} else {
		if len(fields) > 1 {
			if strings.HasSuffix(fields[0], "]") {
				//Root element is an array
				arrayIndexPart := getArrayIndexPart(fields[0])
				fields[0] = arrayIndexPart
				return field.NewMappingField(mapField.HasSepcialField(), mapField.HasArray(), fields), nil
			} else {
				return field.NewMappingField(mapField.HasSepcialField(), mapField.HasArray(), mapField.Getfields()[1:]), nil
			}
		} else {
			//Only attribute name no field name
			return field.NewMappingField(mapField.HasSepcialField(), mapField.HasArray(), []string{}), nil
		}
	}
}

func (m *MappingRef) GetFieldName(mapfield *field.MappingField) (string, error) {
	if mapfield.HasSepcialField() {
		fields := mapfield.Getfields()
		activityNameRef := fields[0]
		if strings.HasPrefix(activityNameRef, "$") {
			return getFieldName(fields[1]), nil
		}
		return getFieldName(fields[0]), nil
	}

	if strings.HasPrefix(m.ref, "$") || strings.Index(m.ref, ".") > 0 {
		log.Debugf("Mapping ref %s", m.ref)
		mappingFields := strings.Split(m.ref, ".")
		if strings.HasPrefix(m.ref, "$") {
			return getFieldName(mappingFields[1]), nil

		}
		log.Debugf("Field name now is: %s", mappingFields[0])
		return getFieldName(mappingFields[0]), nil

	}
	return getFieldName(m.ref), nil
}

func (m *MappingRef) GetActivityId() (string, error) {

	dotIndex := strings.Index(m.ref, ".")

	if dotIndex == -1 {
		return "", fmt.Errorf("invalid resolution expression [%s]", m.ref)
	}

	firstItemIndex := strings.Index(m.ref[:dotIndex], "[")

	if firstItemIndex != -1 {
		return m.ref[firstItemIndex+1 : dotIndex-1], nil
	}
	return "", nil
}

//
//func GetResolutionDetails(toResolve string) (*string, error) {
//
//
//	dotIdx := strings.Index(toResolve, ".")
//
//	if dotIdx == -1 {
//		return nil, fmt.Errorf("invalid resolution expression [%s]", toResolve)
//	}
//
//	details := &ResolutionDetails{}
//	itemIdx := strings.Index(toResolve[:dotIdx], "[")
//
//	if itemIdx != -1 {
//		details.Item = toResolve[itemIdx+1:dotIdx-1]
//		details.ResolverName = toResolve[:itemIdx]
//	} else {
//		details.ResolverName = toResolve[:dotIdx]
//
//		//special case for activity without brackets
//		if strings.HasPrefix(toResolve, "activity") {
//			nextDot := strings.Index(toResolve[dotIdx+1:], ".") + dotIdx + 1
//			details.Item = toResolve[dotIdx+1:nextDot]
//			dotIdx = nextDot
//		}
//	}
//
//	pathIdx := strings.IndexFunc(toResolve[dotIdx+1:], isSep)
//
//	if pathIdx != -1 {
//		pathStart := pathIdx + dotIdx + 1
//		details.Path = toResolve[pathStart:]
//		details.Property = toResolve[dotIdx+1:pathStart]
//	} else {
//		details.Property = toResolve[dotIdx+1:]
//	}
//
//	return details, nil
//}

func getFieldName(fieldname string) string {
	if strings.Index(fieldname, "[") > 0 && strings.Index(fieldname, "]") > 0 {
		return fieldname[:strings.Index(fieldname, "[")]
	}
	return fieldname
}

func getActivityName(fieldname string) string {
	//$activity[name]
	startIndex := strings.Index(fieldname, "[")
	endIndex := strings.Index(fieldname, "]")
	if startIndex >= 0 {
		return fieldname[startIndex+1 : endIndex]
	} else {
		return fieldname
	}
}

//getArrayIndexPart get array part of the string. such as name[0] return [0]
func getArrayIndexPart(fieldName string) string {
	if strings.Index(fieldName, "[") >= 0 {
		return fieldName[strings.Index(fieldName, "[") : strings.Index(fieldName, "]")+1]
	}
	return ""
}
