package data

import (
	"strings"
)

// GetAttrPath splits the supplied attribute with path to its name and object path
func GetAttrPath(inAttrName string) (attrName string, attrPath string) {

	//todo handle bad attr names
	//fmt.Printf("** InAttrName: %s \n", inAttrName)

	if inAttrName[0] == '[' {

		idx := strings.Index(inAttrName, "]")

		if idx == len(inAttrName)-1 {
			attrName = inAttrName
		} else {
			attrName = inAttrName[:idx+1]
			attrPath = inAttrName[idx+2:]
		}
	} else {
		idx := strings.Index(inAttrName, ".")

		if idx == -1 {
			attrName = inAttrName
		} else {
			attrName = inAttrName[:idx]
			attrPath = inAttrName[idx+1:]
		}
	}

	//fmt.Printf("AttrName: %s \n", attrName)
	//fmt.Printf("AttrPath: %s \n", attrPath)

	return attrName, attrPath
}
