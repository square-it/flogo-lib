package random

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/function"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("random-function")

type Random struct {
}

func init() {
	function.Registry(&Random{})
}

func (s *Random) GetName() string {
	return "random"
}

func (s *Random) GetCategory() string {
	return "number"
}

func (s *Random) Eval(limitIn interface{}) int64 {
	limit, err := ConvertToInt64(limitIn)
	if err != nil {
		log.Errorf("Convert %+v to int error %s", limitIn, err.Error())
		limit = 10
	}
	log.Debugf("Generate sudo-random integer number within the scope of [0, %d)", limit)
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(limit)
}

func ConvertToInt64(value interface{}) (int64, error) {
	switch t := value.(type) {
	case string:
		return strconv.ParseInt(t, 10, 64)
	case int:
		return int64(t), nil
	case int64:
		return t, nil
	case float64:
		return int64(t), nil
	default:
		v, err := json.Marshal(value)
		if err != nil {
			return 0, err
		}
		i, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return 0, err
		}
		return i, nil
	}
}
