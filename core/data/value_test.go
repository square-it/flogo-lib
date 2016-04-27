package data

import "testing"

func TestGetAttrPath(t *testing.T) {

	a:= "sensorData.temp"
	GetAttrPath(a)

	a= "T.v"
	GetAttrPath(a)

}
