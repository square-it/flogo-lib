package link

import (
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/stretchr/testify/assert"
)

func TestWILinkExprManager_EvalLinkExpr(t *testing.T) {
	link := &definition.Link{}
	link.Value()
	factory := &WILinkExprManagerFactory{}
	b, err := factory.NewLinkExprManager(nil).EvalLinkExpr(link, nil)
	assert.Nil(t, err)
	assert.True(t, b)
}
