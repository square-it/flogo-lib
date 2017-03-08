package data

import (
	"strings"
)

// PathType is the attribute value accessor path
type PathType int

const (
	PT_SIMPLE PathType = 1
	PT_MAP    PathType = 2
	PT_ARRAY  PathType = 3
)

// GetAttrPath splits the supplied attribute with path to its name and object path
func GetAttrPath(inAttrName string) (attrName string, attrPath string, pathType PathType) {

	//todo handle bad attr names
	//fmt.Printf("** InAttrName: %s \n", inAttrName)

	nameLen := len(inAttrName)
	pathType = PT_SIMPLE

	if inAttrName[0] == '{' {

		idx := strings.Index(inAttrName, "}")

		if idx == nameLen-1 {
			attrName = inAttrName
		} else {
			attrName = inAttrName[:idx+1]

			if inAttrName[idx+1] == '[' {
				pathType = PT_ARRAY
				attrPath = inAttrName[idx+2 : nameLen-1]
			} else {
				pathType = PT_MAP
				attrPath = inAttrName[idx+2:]
			}
		}
	} else {
		idx := strings.Index(inAttrName, ".")

		if idx == -1 {

			idx = strings.Index(inAttrName, "[")

			if idx == -1 {
				attrName = inAttrName
			} else {
				pathType = PT_ARRAY
				attrName = inAttrName[:idx]
				attrPath = inAttrName[idx+1 : nameLen-1]
			}
		} else {
			pathType = PT_MAP
			attrName = inAttrName[:idx]
			attrPath = inAttrName[idx+1:]
		}
	}

	return attrName, attrPath, pathType
}

func GetMapValue(valueMap map[string]interface{}, path string) interface{} {

	var pathComponents []string = strings.Split(path, ".")
	lastPcIdx := len(pathComponents) - 1

	tmpObj := valueMap

	for pcIdx, pc := range pathComponents {

		if pcIdx == lastPcIdx {
			return tmpObj[pc]
		}

		switch tmpObj[pc].(type) {
		//todo need to throw error if not a map
		case map[string]interface{}:
			tmpObj = tmpObj[pc].(map[string]interface{})
		}
	}

	return tmpObj
}
