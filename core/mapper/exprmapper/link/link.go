package link

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("LinkExpr")

type WILinkExprManager struct {
}

type WILinkExprManagerFactory struct {
}

// NewGosLinkExprManager creates a new LuaLinkExprManager
func (f *WILinkExprManagerFactory) NewLinkExprManager(def *definition.Definition) definition.LinkExprManager {
	mgr := &WILinkExprManager{}
	return mgr
}

// EvalLinkExpr implements LinkExprManager.EvalLinkExpr
func (em *WILinkExprManager) EvalLinkExpr(link *definition.Link, scope data.Scope) (bool, error) {
	value := link.Value()
	if value == "" {
		return true, nil
	}

	log.Debugf("WI link expression value [%s]", value)
	//Todo inject resolver
	funcValue, err := exprmapper.GetMappingValue(value, scope, nil)
	if err != nil {
		log.Error("Get value from link value %+v, error %s", value, err.Error())
		return false, fmt.Errorf("Get value from link value %+v, error %s", value, err.Error())
	}

	b, err := data.CoerceToBoolean(funcValue)
	if err != nil {
		log.Error("Parser [%+v] to boolean error [%s]", value, err.Error())
		return false, fmt.Errorf("Parser [%+v] to boolean error [%s]", value, err.Error())
	}
	log.Debugf("Linking %s result %b", link.Value(), b)
	return b, nil
}
