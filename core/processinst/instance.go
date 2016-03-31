package processinst

import (
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/ext/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/ext/model"
	"github.com/TIBCOSoftware/flogo-lib/core/process"
	"github.com/TIBCOSoftware/flogo-lib/util"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("instance")

// Instance is a structure for representing an instance of a Process
type Instance struct {
	id           string
	stepID       int
	lock         sync.Mutex
	status       Status
	state        int
	ProcessURI   string
	Process      *process.Definition
	RootTaskEnv  *TaskEnv
	ProcessModel *model.ProcessModel
	Attrs        map[string]*data.Attribute
	Patch        *process.Patch
	Interceptor  *process.Interceptor

	WorkItemQueue *util.SyncQueue //todo: change to faster non-threadsafe queue

	wiCounter     int
	ChangeTracker *InstanceChangeTracker `json:"-"`

	processProvider process.Provider
}

// NewProcessInstance creates a new Process Instance from the specified Process
func NewProcessInstance(instanceID string, processURI string, process *process.Definition) *Instance {

	var instance Instance
	instance.id = instanceID
	instance.stepID = 0
	instance.ProcessURI = processURI
	instance.Process = process
	instance.ProcessModel = model.Get(process.ModelID())
	instance.status = StatusNotStarted
	instance.WorkItemQueue = util.NewQueue()
	instance.ChangeTracker = NewInstanceChangeTracker()

	var taskEnv TaskEnv
	taskEnv.ID = 1
	taskEnv.Task = process.RootTask()
	taskEnv.taskID = process.RootTask().ID()
	taskEnv.Instance = &instance
	taskEnv.TaskDatas = make(map[int]*TaskData)
	taskEnv.LinkDatas = make(map[int]*LinkData)

	instance.RootTaskEnv = &taskEnv

	return &instance
}

func (pi *Instance) SetProcessProvidcer(provider process.Provider) {
	pi.processProvider = provider
}

// Restart indicates that this ProcessInstance was restarted
func (pi *Instance) Restart(id string, provider process.Provider) {
	pi.id = id
	pi.processProvider = provider
	pi.Process = pi.processProvider.GetProcess(pi.ProcessURI)
	pi.ProcessModel = model.Get(pi.Process.ModelID())
	pi.RootTaskEnv.init(pi)
}

// ID returns the ID of the Process Instance
func (pi *Instance) ID() string {
	return pi.id
}

// ProcessDefinition returns the Process that the instance is of
func (pi *Instance) ProcessDefinition() *process.Definition {
	return pi.Process
}

// StepID returns the current step ID of the Process Instance
func (pi *Instance) StepID() int {
	return pi.stepID
}

// Status returns the current status of the Process Instance
func (pi *Instance) Status() Status {
	return pi.status
}

func (pi *Instance) setStatus(status Status) {

	pi.status = status
	pi.ChangeTracker.SetStatus(status)
}

// State returns the state indicator of the Process Instance
func (pi *Instance) State() int {
	return pi.state
}

// SetState sets the state indicator of the Process Instance
func (pi *Instance) SetState(state int) {
	pi.state = state
	pi.ChangeTracker.SetState(state)
}

// UpdateAttrs updates the attributes of the Process Instance
func (pi *Instance) UpdateAttrs(update interface{}) {

	if update != nil {

		log.Debugf("Updating process attrs: %v", update)

		m := update.(map[string]string)

		if pi.Attrs == nil {
			pi.Attrs = make(map[string]*data.Attribute, len(m))
		}

		for k, v := range m {
			pi.Attrs[k] = &data.Attribute{Name: k, Type: "string", Value: v}
		}
	}
}

// Start will start the Process Instance, returns a boolean indicating
// if it was able to start
func (pi *Instance) Start(data interface{}) bool {

	pi.setStatus(StatusActive)
	pi.UpdateAttrs(data)

	log.Infof("ProcessInstance Process: %v", pi.ProcessModel)
	model := pi.ProcessModel.GetProcessBehavior(pi.Process.TypeID())

	ok, evalCode := model.Start(pi, data)

	if ok {
		rootTaskData := pi.RootTaskEnv.NewTaskData(pi.Process.RootTask())

		pi.scheduleEval(rootTaskData, evalCode)
	}

	return ok
}

//Resume resumes a Process Instance
func (pi *Instance) Resume(data interface{}) bool {

	model := pi.ProcessModel.GetProcessBehavior(pi.Process.TypeID())

	return model.Resume(pi, data)
}

// DoStep performs a single execution 'step' of the Process Instance
func (pi *Instance) DoStep() bool {

	hasNext := false

	pi.ResetChanges()

	pi.stepID++

	if pi.status == StatusActive {

		item, ok := pi.WorkItemQueue.Pop()

		if ok {
			log.Debug("popped item off queue")

			workItem := item.(*WorkItem)

			pi.ChangeTracker.trackWorkItem(&WorkItemQueueChange{ChgType: CtDel, ID: workItem.ID, WorkItem: workItem})

			pi.execTask(workItem)
			hasNext = true
		} else {
			log.Debug("queue emtpy")
		}
	}

	return hasNext
}

// GetChanges returns the Change Tracker object
func (pi *Instance) GetChanges() *InstanceChangeTracker {
	return pi.ChangeTracker
}

// ResetChanges resets an changes that were being tracked
func (pi *Instance) ResetChanges() {

	if pi.ChangeTracker != nil {
		pi.ChangeTracker.ResetChanges()
	}

	//todo: can we reuse this to avoid gc
	pi.ChangeTracker = NewInstanceChangeTracker()
}

func (pi *Instance) scheduleEval(taskData *TaskData, evalCode int) {

	pi.wiCounter++

	workItem := NewWorkItem(pi.wiCounter, taskData, EtEval, evalCode)
	log.Debugf("Scheduling EVAL on task: %s\n", taskData.task.Name())

	pi.WorkItemQueue.Push(workItem)
	pi.ChangeTracker.trackWorkItem(&WorkItemQueueChange{ChgType: CtAdd, ID: workItem.ID, WorkItem: workItem})
}

// execTask executes the specified Work Item of the Process Instance
func (pi *Instance) execTask(workItem *WorkItem) {
	
	taskBehavior := pi.ProcessModel.GetTaskBehavior(workItem.TaskData.task.TypeID())

	var done bool
	var doneCode int

	if workItem.ExecType == EtEval {

		// get the input mapper
		inputMapper := workItem.TaskData.task.InputMapper()

		if pi.Patch != nil {
			// check if the patch has a overriding mapper
			mapper := pi.Patch.GetInputMapper(workItem.TaskData.task.ID())
			if mapper != nil {
				inputMapper = mapper
			}
		}

		if inputMapper != nil {
			log.Debug("Applying InputMapper")
			inputMapper.Apply(pi, workItem.TaskData)
		}

		eval := true

		if pi.Interceptor != nil {
			// check if this task as an interceptor
			taskInterceptor := pi.Interceptor.GetTaskInterceptor(workItem.TaskData.task.ID())

			if taskInterceptor != nil {

				log.Debug("Applying Interceptor")

				if len(taskInterceptor.Inputs) > 0 {
					// override input attributes
					for _, attribute := range taskInterceptor.Inputs {

						log.Debugf("Overriding Attr: %s = %s", attribute.Name, attribute.Value)

						//todo: validation
						workItem.TaskData.SetAttrValue(attribute.Name, attribute.Value)
					}
				}

				// check if we should not evaluate the task
				eval = !taskInterceptor.Skip
			}
		}

		if eval {
			done, doneCode = taskBehavior.Eval(workItem.TaskData, workItem.EvalCode)
		} else {
			done = true
		}
	} else {
		done, doneCode = taskBehavior.PostEval(workItem.TaskData, workItem.EvalCode, nil)
	}

	if done {

		if pi.Interceptor != nil {
			// check if this task as an interceptor and overrides ouputs
			taskInterceptor := pi.Interceptor.GetTaskInterceptor(workItem.TaskData.task.ID())
			if taskInterceptor != nil && len(taskInterceptor.Outputs) > 0 {
				// override output attributes
				for _, attribute := range taskInterceptor.Outputs {

					//todo: validation
					workItem.TaskData.SetAttrValue(attribute.Name, attribute.Value)
				}
			}
		}

		// get the Output Mapper for the Task if one exists
		outputMapper := workItem.TaskData.task.OutputMapper()

		if pi.Patch != nil {
			// check if the patch overrides the Output Mapper
			mapper := pi.Patch.GetOutputMapper(workItem.TaskData.task.ID())
			if mapper != nil {
				outputMapper = mapper
			}
		}

		if outputMapper != nil {
			log.Debug("Applying OutputMapper")
			outputMapper.Apply(workItem.TaskData, pi)
		}

		pi.handleTaskDone(taskBehavior, workItem.TaskData, doneCode)
	}
}

// handleTaskDone handles the compeletion of a task in the Process Instance
func (pi *Instance) handleTaskDone(taskBehavior model.TaskBehavior, taskData *TaskData, doneCode int) {

	notifyParent, childDoneCode, taskEntries := taskBehavior.Done(taskData, doneCode)

	task := taskData.Task()

	if notifyParent {

		parentTask := task.Parent()

		if parentTask != nil {
			parentTaskData := taskData.taskEnv.TaskDatas[parentTask.ID()]
			parentBehavior := pi.ProcessModel.GetTaskBehavior(parentTask.TypeID())
			parentDone, parentDoneCode := parentBehavior.ChildDone(parentTaskData, task, childDoneCode)

			if parentDone {
				pi.handleTaskDone(parentBehavior, parentTaskData, parentDoneCode)
			}

		} else {
			//Root Task is Done, so notify Process
			processBehavior := pi.ProcessModel.GetProcessBehavior(pi.Process.TypeID())
			processBehavior.TasksDone(pi, childDoneCode)
			processBehavior.Done(pi)

			pi.setStatus(StatusCompleted)
		}
	}

	if len(taskEntries) > 0 {

		for _, taskEntry := range taskEntries {

			log.Debugf("execTask - TaskEntry: %v\n", taskEntry)
			taskToEnterBehavior := pi.ProcessModel.GetTaskBehavior(taskEntry.Task.TypeID())

			enterTaskData, _ := taskData.taskEnv.FindOrCreateTaskData(taskEntry.Task)

			eval, evalCode := taskToEnterBehavior.Enter(enterTaskData, taskEntry.EnterCode)

			if eval {
				pi.scheduleEval(enterTaskData, evalCode)
			}
		}
	}

	taskData.taskEnv.releaseTask(task)
}

// GetAttrType implements api.Scope.GetAttrType
func (pi *Instance) GetAttrType(attrName string) (attrType string, exists bool) {

	if pi.Attrs != nil {
		attr, found := pi.Attrs[attrName]
		if found {
			return attr.Type, true
		}
	}

	attr, found := pi.Process.GetAttr(attrName)
	if found {
		return attr.Type, true
	}

	return "", false
}

// GetAttrValue implements api.Scope.GetAttrValue
func (pi *Instance) GetAttrValue(attrName string) (value string, exists bool) {

	if pi.Attrs != nil {
		attr, found := pi.Attrs[attrName]
		if found {
			return attr.Value, true
		}
	}

	attr, found := pi.Process.GetAttr(attrName)
	if found {
		return attr.Value, true
	}

	return "", false
}

// SetAttrValue implements api.Scope.SetAttrValue
func (pi *Instance) SetAttrValue(attrName string, value string) {
	if pi.Attrs == nil {
		pi.Attrs = make(map[string]*data.Attribute)
	}

	attrType, exists := pi.GetAttrType(attrName)

	if exists {
		pi.Attrs[attrName] = &data.Attribute{Name: attrName, Type: attrType, Value: value}
	}
	// else what do we do if its a completely new attr
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// Task Environment

// TaskEnv is a structure that describes the execution enviroment for a set of tasks
type TaskEnv struct {
	ID        int
	Task      *process.Task
	Instance  *Instance
	ParentEnv *TaskEnv

	TaskDatas map[int]*TaskData
	LinkDatas map[int]*LinkData

	taskID int // for deserialization
}

// init initializes the Task Environment, typically called on deserialization
func (te *TaskEnv) init(processInst *Instance) {

	if te.Instance == nil {

		te.Instance = processInst
		te.Task = processInst.Process.GetTask(te.taskID)

		for _, v := range te.TaskDatas {
			v.taskEnv = te
			v.task = processInst.Process.GetTask(v.taskID)
		}

		for _, v := range te.LinkDatas {
			v.taskEnv = te
			v.link = processInst.Process.GetLink(v.linkID)
		}
	}
}

// FindOrCreateTaskData finds an existing TaskData or creates ones if not found for the
// specified task the task enviroment
func (te *TaskEnv) FindOrCreateTaskData(task *process.Task) (taskData *TaskData, created bool) {

	taskData, ok := te.TaskDatas[task.ID()]

	created = false

	if !ok {
		taskData = NewTaskData(te, task)
		te.TaskDatas[task.ID()] = taskData
		te.Instance.ChangeTracker.trackTaskData(&TaskDataChange{ChgType: CtAdd, ID: task.ID(), TaskData: taskData})

		created = true
	}

	return taskData, created
}

// NewTaskData creates a new TaskData object
func (te *TaskEnv) NewTaskData(task *process.Task) *TaskData {

	taskData := NewTaskData(te, task)
	te.TaskDatas[task.ID()] = taskData
	te.Instance.ChangeTracker.trackTaskData(&TaskDataChange{ChgType: CtAdd, ID: task.ID(), TaskData: taskData})

	return taskData
}

// FindOrCreateLinkData finds an existing LinkData or creates ones if not found for the
// specified link the task enviroment
func (te *TaskEnv) FindOrCreateLinkData(link *process.Link) (linkData *LinkData, created bool) {

	linkData, ok := te.LinkDatas[link.ID()]
	created = false

	if !ok {
		linkData = NewLinkData(te, link)
		te.LinkDatas[link.ID()] = linkData
		te.Instance.ChangeTracker.trackLinkData(&LinkDataChange{ChgType: CtAdd, ID: link.ID(), LinkData: linkData})
		created = true
	}

	return linkData, created
}

// releaseTask cleans up TaskData in the task environment any of its dependencies.
// This is called when a task is completed and can be discarded
func (te *TaskEnv) releaseTask(task *process.Task) {
	delete(te.TaskDatas, task.ID())
	te.Instance.ChangeTracker.trackTaskData(&TaskDataChange{ChgType: CtDel, ID: task.ID()})

	childTasks := task.ChildTasks()

	if len(childTasks) > 0 {

		for _, childTask := range childTasks {
			delete(te.TaskDatas, childTask.ID())
			te.Instance.ChangeTracker.trackTaskData(&TaskDataChange{ChgType: CtDel, ID: childTask.ID()})
		}
	}

	links := task.FromLinks()

	for _, link := range links {
		delete(te.LinkDatas, link.ID())
		te.Instance.ChangeTracker.trackLinkData(&LinkDataChange{ChgType: CtDel, ID: link.ID()})
	}
}

// TaskData represents data associated with an instance of a Task
type TaskData struct {
	taskEnv *TaskEnv
	task    *process.Task
	state   int
	done    bool
	attrs   map[string]*data.Attribute

	changes int

	taskID int //needed for serialization
}

// NewTaskData creates a TaskData for the specified task in the specified task
// environment
func NewTaskData(taskEnv *TaskEnv, task *process.Task) *TaskData {
	var taskData TaskData

	taskData.taskEnv = taskEnv
	taskData.task = task
	//taskData.TaskID = task.ID

	return &taskData
}

/////////////////////////////////////////
// TaskData - TaskContext Implementation

// State implements process.TaskContext.GetState
func (td *TaskData) State() int {
	return td.state
}

// SetState implements process.TaskContext.SetState
func (td *TaskData) SetState(state int) {
	td.state = state
	td.taskEnv.Instance.ChangeTracker.trackTaskData(&TaskDataChange{ChgType: CtUpd, ID: td.task.ID(), TaskData: td})
}

// Task implements model.TaskContext.Task, by returning the Task associated with this
// TaskData object
func (td *TaskData) Task() *process.Task {
	return td.task
}

// FromLinks implements model.TaskContext.GetFromLinks, by returning the set of predecessor
// Links of the current task.
func (td *TaskData) FromLinks() []model.LinkContext {

	log.Debugf("GetFromLinks: task=%v\n", td.Task)

	links := td.task.FromLinks()

	numLinks := len(links)

	if numLinks > 0 {
		linkCtxs := make([]model.LinkContext, numLinks)

		for i, link := range links {
			linkCtxs[i], _ = td.taskEnv.FindOrCreateLinkData(link)
		}
		return linkCtxs
	}

	return nil
}

// EnterChildren implements activity.ActivityContext.EnterChildren method
func (td *TaskData) EnterChildren(taskEntries []*model.TaskEntry) {

	if (taskEntries == nil) || (len(taskEntries) == 1 && taskEntries[0].Task == nil) {

		var enterCode int

		if taskEntries == nil {
			enterCode = 0
		} else {
			enterCode = taskEntries[0].EnterCode
		}

		log.Debugf("Entering '%s' Task's %d children\n", td.task.Name(), len(td.task.ChildTasks()))

		for _, task := range td.task.ChildTasks() {

			taskData, _ := td.taskEnv.FindOrCreateTaskData(task)
			taskBehavior := td.taskEnv.Instance.ProcessModel.GetTaskBehavior(task.TypeID())

			eval, evalCode := taskBehavior.Enter(taskData, enterCode)

			if eval {
				td.taskEnv.Instance.scheduleEval(taskData, evalCode)
			}
		}
	} else {

		for _, taskEntry := range taskEntries {

			//todo validate if specified task is child? or trust model

			taskData, _ := td.taskEnv.FindOrCreateTaskData(taskEntry.Task)
			taskBehavior := td.taskEnv.Instance.ProcessModel.GetTaskBehavior(taskEntry.Task.TypeID())

			eval, evalCode := taskBehavior.Enter(taskData, taskEntry.EnterCode)

			if eval {
				td.taskEnv.Instance.scheduleEval(taskData, evalCode)
			}
		}
	}
}

// EvalLink implements activity.ActivityContext.EvalLink method
func (td *TaskData) EvalLink(link *process.Link, evalCode int) model.LinkContext {

	linkData, _ := td.taskEnv.FindOrCreateLinkData(link)

	//linkBehavior := td.taskEnv.Instance.ProcessModel.GetLinkBehavior(link.)
	//linkBehavior.Eval(linkData, evalCode)
	linkData.SetState(2)

	log.Debugf("TaskContext.EvalLink: State = %d\n", linkData.State())

	return linkData
}

// Activity implements activity.Context.Activity method
func (td *TaskData) Activity() (act activity.Activity, context activity.Context) {

	act = activity.Get(td.task.ActivityType())

	return act, td
}

// ProcessInstanceID implements activity.Context.ProcessInstanceID method
func (td *TaskData) ProcessInstanceID() string {

	return td.taskEnv.Instance.id
}

// ProcessName implements activity.Context.ProcessName method
func (td *TaskData) ProcessName() string {
	return td.taskEnv.Instance.Process.Name()
}

// TaskName implements activity.Context.TaskName method
func (td *TaskData) TaskName() string {
	return td.task.Name()
}

// GetAttrType implements api.Scope.GetAttrType
func (td *TaskData) GetAttrType(attrName string) (attrType string, exists bool) {

	if td.attrs != nil {
		attr, found := td.attrs[attrName]
		if found {
			return attr.Type, true
		}
	}

	attr, found := td.task.GetAttr(attrName)
	if found {
		return attr.Type, true
	}

	return "", false
}

// GetAttrValue implements api.Scope.GetAttrValue
func (td *TaskData) GetAttrValue(attrName string) (value string, exists bool) {
	if td.attrs != nil {
		attr, found := td.attrs[attrName]
		if found {
			return attr.Value, true
		}
	}

	attr, found := td.task.GetAttr(attrName)
	if found {
		return attr.Value, true
	}

	return "", false
}

// SetAttrValue implements api.Scope.SetAttrValue
func (td *TaskData) SetAttrValue(attrName string, value string) {

	if td.attrs == nil {
		td.attrs = make(map[string]*data.Attribute)
	}

	attrType, exists := td.GetAttrType(attrName)

	if exists {
		td.attrs[attrName] = &data.Attribute{Name: attrName, Type: attrType, Value: value}
	}
	// todo: else what do we do if its a completely new attr, how should we infer the type
}

// LinkData represents data associated with an instance of a Link
type LinkData struct {
	taskEnv *TaskEnv
	link    *process.Link
	state   int
	attrs   map[string]string

	changes int

	linkID int //needed for serialization
}

// NewLinkData creates a LinkData for the specified link in the specified task
// environment
func NewLinkData(taskEnv *TaskEnv, link *process.Link) *LinkData {
	var linkData LinkData

	linkData.taskEnv = taskEnv
	linkData.link = link
	//linkData.LinkID = link.ID

	return &linkData
}

// State returns the current state indicator for the LinkData
func (ld *LinkData) State() int {
	return ld.state
}

// SetState sets the current state indicator for the LinkData
func (ld *LinkData) SetState(state int) {
	ld.state = state
	ld.taskEnv.Instance.ChangeTracker.trackLinkData(&LinkDataChange{ChgType: CtUpd, ID: ld.link.ID(), LinkData: ld})
}

// Link returns the Link associated with ld context
func (ld *LinkData) Link() *process.Link {

	return ld.link
}

// ExecType is the type of execution to perform
type ExecType int

const (
	// EtEval denoted the Eval execution type
	EtEval ExecType = 10

	// EtPostEval denoted the PostEval execution type
	EtPostEval ExecType = 20
)

// WorkItem describes an item of work (event for a Task) that should be executed on Step
type WorkItem struct {
	ID       int       `json:"id"`
	TaskData *TaskData `json:"-"`
	ExecType ExecType  `json:"execType"`
	EvalCode int       `json:"code"`

	TaskID int `json:"taskID"` //for now need for ser
	//taskCtxID int `json:"taskCtxID"` //not needed for now
}

// NewWorkItem constructs a new WorkItem for the specified TaskData
func NewWorkItem(id int, taskData *TaskData, execType ExecType, evalCode int) *WorkItem {

	var workItem WorkItem

	workItem.ID = id
	workItem.TaskData = taskData
	workItem.ExecType = execType
	workItem.EvalCode = evalCode

	workItem.TaskID = taskData.task.ID()

	return &workItem
}
