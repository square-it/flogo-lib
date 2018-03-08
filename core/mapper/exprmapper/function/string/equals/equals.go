package equals

import (
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/function"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("equals-function")

type Equals struct {
}

func init() {
	function.Registry(&Equals{})
}

func (s *Equals) GetName() string {
	return "equals"
}

func (s *Equals) GetCategory() string {
	return "string"
}
func (s *Equals) Eval(str, str2 string, ignoreCase bool) bool {
	log.Debugf(`Reports whether "%s" equels "%s" with ignore case %s`, str, str2, ignoreCase)
	if ignoreCase {
		return strings.EqualFold(str, str2)
	}
	return str == str2
}
