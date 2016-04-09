package processinst

// StateRecorder is the interface that describes a service that can record
// snapshots and steps of a Process Instance
type StateRecorder interface {

	// RecordSnapshot records a Snapshot of the ProcessInstance
	RecordSnapshot(instance *Instance)

	// RecordStep records the changes for the current Step of the Process Instance
	RecordStep(instance *Instance)
}
