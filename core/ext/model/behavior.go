package model

import (
	"github.com/TIBCOSoftware/flogo-lib/core/process"
)

// TaskEntry is a struct used to specify what Task to
// enter and its corresponding enter code
type TaskEntry struct {
	Task      *process.Task
	EnterCode int
}

// ProcessBehavior is the execution behavior of the Process.
type ProcessBehavior interface {

	// Start the process instance.  Returning true indicates that the
	// process can start and eval will be scheduled on the Root Task.
	// Return false indicates that the process could not be started
	// at this time.
	Start(context ProcessContext, data interface{}) (start bool, evalCode int)

	// Resume the process instance.  Returning true indicates that the
	// process can resume.  Return false indicates that the process
	// could not be resumed at this time.
	Resume(context ProcessContext, data interface{}) bool //<---

	//do we need the following two

	// TasksDone is called when the RootTask is Done.
	TasksDone(context ProcessContext, doneCode int)

	// Done is called when the process is done.
	Done(context ProcessContext) //maybe return something to the state server?
}

// TaskBehavior is the execution behavior of a Task.
type TaskBehavior interface {

	// Enter determines if a Task is ready to be evaluated, returning true
	// indicates that the task is ready to be evaluated.
	Enter(context TaskContext, enterCode int) (eval bool, evalCode int)

	// Eval is called when a Task is being evaluated.  Returning true indicates
	// that the task is done.
	Eval(context TaskContext, evalCode int) (done bool, doneCode int)

	// PostEval is called when a task that didn't complete during the Eval
	// needs to be notified.  Returning true indicates that the task is done.
	PostEval(context TaskContext, evalCode int, data interface{}) (done bool, doneCode int)

	// Done is called when Eval, PostEval or ChildDone return true, indicating
	// that the task is done.  This step is used to finalize the task and
	// determine the next set of tasks to be entered.  Returning true indicates
	// that the parent task should be notified.  Also returns the set of Tasks
	// that should be entered next.
	Done(context TaskContext, doneCode int) (notifyParent bool, childDoneCode int, taskEntries []*TaskEntry)

	// ChildDone is called when child task is Done and has indicated that its
	// parent should be notified.  Returning true indicates that the task
	// is done.
	ChildDone(context TaskContext, childTask *process.Task, childDoneCode int) (done bool, doneCode int)
}

// LinkBehavior is the execution behavior of a Link.
// todo: remove, link do not need behaviors
type LinkBehavior interface {

	// Eval is called when a Link is being evaluated.
	Eval(context LinkContext, evalCode int)
}
