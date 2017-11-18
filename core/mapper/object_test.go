package mapper

import (
	"testing"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

func Test_evalExpr(t *testing.T) {
	evalExpr("{{1}}", nil, &data.BasicResolver{})
}