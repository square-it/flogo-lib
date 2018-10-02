package json

import (
	"fmt"
	"testing"
)

type TestStruct struct {
	Id      string `json:"id"`
	name    string `json:"name"`
	Zipcode string `json:"zipcode"`
	Test    string `json:",omitempty`
}

func TestGetFieldByName(t *testing.T) {
	test := &TestStruct{Id: "idssss", name: "dddddd", Zipcode: "55555", Test: "ddd"}
	field, _ := GetFieldByName(test, "id")
	fmt.Println(field.Interface())
}
