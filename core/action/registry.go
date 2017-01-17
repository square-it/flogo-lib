package action

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/op/go-logging"
)

var (
	actionsMu sync.Mutex
	actions   = make(map[string]Action)
	log       = logging.MustGetLogger("action")
	reg       = &registry{}
)

type Registry interface {
	Add(t Action2) error
	GetActions() map[string]Action2
	GetAction(string) Action2
}

type registry struct {
	actions map[string]Action2
}

func GetRegistry() Registry {
	return reg
}

func (r *registry) Add(a Action2) error {
	actionsMu.Lock()
	defer actionsMu.Unlock()

	if a == nil {
		return fmt.Errorf("registry.Add: trigger is nil")
	}

	// copy on write to avoid synchronization on access
	newAs := make(map[string]Action2, len(r.actions))

	for k, v := range r.actions {
		newAs[k] = v
	}

	AddAction(newAs, a)
	r.actions = newAs

	return nil
}

func AddAction(m map[string]Action2, action Action2) {
	t := reflect.TypeOf(action)
	pkgPath := t.Elem().PkgPath()
	pkgPath = strings.TrimLeft(pkgPath, "vendor/src/")
	pkgPath = strings.TrimLeft(pkgPath, "vendor/")
	m[pkgPath] = action
}

//GetActions returns a map of all the registered Actions where key is the pkg name of the type
func (r *registry) GetActions() map[string]Action2 {
	return r.actions
}

// Register registers the specified action
func Register(actionType string, action Action) {
	actionsMu.Lock()
	defer actionsMu.Unlock()

	if actionType == "" {
		panic("action.Register: actionType is empty")
	}

	if action == nil {
		panic("action.Register: action is nil")
	}

	if _, dup := actions[actionType]; dup {
		panic("action.Register: action already registered for action type: " + actionType)
	}

	// copy on write to avoid synchronization on access
	newActions := make(map[string]Action, len(actions))

	for k, v := range actions {
		newActions[k] = v
	}

	newActions[actionType] = action
	actions = newActions

	log.Debugf("Registerd Action: %s", actionType)
}

// Actions gets all the registered Action Actions
func Actions() []Action {

	var curActions = actions

	list := make([]Action, 0, len(curActions))

	for _, value := range curActions {
		list = append(list, value)
	}

	return list
}

// Get gets specified Action
func Get(actionType string) Action {
	return actions[actionType]
}

// Get gets specified Action
func (r *registry) GetAction(id string) Action2 {
	return r.actions[id]
}
