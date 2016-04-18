package activity

// Activity is an interface for defining a custom Task Execution
type Activity interface {

	// Eval is called when an Activity is being evaluated.  Returning true indicates
	// that the task is done.
	Eval(context Context) bool //, Error - suspend, or follow error transition?

	// ActivityMetadata returns the metadata of the activity
	Metadata() *Metadata
}
