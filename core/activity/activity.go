package activity

// Activity is an interface for defining a custom Activity Execution
type Activity interface {

	// Eval is called when an Activity is being evaluated.  Returning true indicates
	// that the task is done.
	Eval(ctx Context) (done bool, err error)

	// ActivityMetadata returns the metadata of the activity
	Metadata() *Metadata
}

// Init is an optional interface that can be implemented by an activity.  If implemented,
// it will be invoked for each corresponding activity configuration that has settings.
type Init interface {

	// Init initialize the Activity for a particular configuration
	Init(ctx InitContext) error
}