package action

import (
	"sync"

	"github.com/op/go-logging"
)

var (
	actionsMu sync.Mutex
	actions   = make(map[string]Action)
)

var log = logging.MustGetLogger("action")

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
