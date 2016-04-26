package flow

import "github.com/TIBCOSoftware/flogo-lib/core/data"

//todo rename
type LinkExprManager interface {

	evalLinkExpr(link *Link, scope data.Scope) bool

}

