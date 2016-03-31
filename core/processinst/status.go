package processinst

// Status is value that indicates the status of a Process Instance
type Status int

const (
	// StatusNotStarted indicates that the ProcessInstance has not started
	StatusNotStarted Status = 0

	// StatusActive indicates that the ProcessInstance is active
	StatusActive Status = 100

	// StatusCompleted indicates that the ProcessInstance has been completed
	StatusCompleted Status = 500

	// StatusCancelled indicates that the ProcessInstance has been cancelled
	StatusCancelled Status = 600

	// StatusFailed indicates that the ProcessInstance has failed
	StatusFailed Status = 700
)
