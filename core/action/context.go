package action

import "github.com/TIBCOSoftware/flogo-lib/core/data"

// Context describes the execution context for an object within an Action.
// It provides access to its configuration and instance information..
type Context interface {

	// ID returns the ID of the Action Instance
	ID() string

	// The action reference
	Ref() string

	// Get metadata of the action instance
	InstanceMetadata() *ConfigMetadata

	// Reply is used to reply with the results of the instance execution
	Reply(replyData map[string]interface{}, err error)

	// todo see which reply makes more sense
	ReplyWithAttrs(replyData map[string]*data.Attribute, err error)

	//todo rename, essentially the flow's attrs for now
	WorkingData() data.Scope

	//todo add a complete method

	//Map with action specific properties
}

