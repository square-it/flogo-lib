package mapper

import "testing"

func Test_evalExpr(t *testing.T) {
	evalExpr("{{1}}", nil)
}