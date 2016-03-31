package process

import (
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("process")

// Definition is the object that describes the definition of
// a process.  It contains its data (attributes) and
// structure (tasks & links).
type Definition struct {
	typeID   int
	name     string
	modelID  string
	rootTask *Task

	attrs map[string]*data.Attribute

	inputMapper *data.Mapper
	links       map[int]*Link
	tasks       map[int]*Task
}

// Name returns the name of the definition
func (pd *Definition) Name() string {
	return pd.name
}

// TypeID returns the type ID of the definition
func (pd *Definition) TypeID() int {
	return pd.typeID
}

// ModelID returns the ID of the model the definition uses
func (pd *Definition) ModelID() string {
	return pd.modelID
}

// RootTask returns the root task of the definition
func (pd *Definition) RootTask() *Task {
	return pd.rootTask
}

// GetAttr gets the specified attribute
func (pd *Definition) GetAttr(attrName string) (attr *data.Attribute, exists bool) {

	if pd.attrs != nil {
		attr, found := pd.attrs[attrName]
		if found {
			return attr, true
		}
	}

	return nil, false
}

// GetTask returns the task with the specified ID
func (pd *Definition) GetTask(taskID int) *Task {
	task := pd.tasks[taskID]
	return task
}

// GetLink returns the link with the specified ID
func (pd *Definition) GetLink(linkID int) *Link {
	task := pd.links[linkID]
	return task
}

////////////////////////////////////////////////////////////////////////////
// Task

// Task is the object that describes the definition of
// a task.  It contains its data (attributes) and its
// nested structure (child tasks & child links).
type Task struct {
	id           int
	typeID       int
	activityType string
	name         string
	tasks        []*Task
	links        []*Link

	definition   *Definition
	parent       *Task
	attrs        map[string]*data.Attribute

	inputMapper  *data.Mapper
	outputMapper *data.Mapper

	toLinks      []*Link
	fromLinks    []*Link
}

func (task *Task) ID() int {
	return task.id
}

func (task *Task) Name() string {
	return task.name
}

func (task *Task) TypeID() int {
	return task.typeID
}

func (task *Task) ActivityType() string {
	return task.activityType
}

func (task *Task) Parent() *Task {
	return task.parent
}

func (task *Task) ChildTasks() []*Task {
	return task.tasks
}

// GetAttr gets the specified attribute
func (task *Task) GetAttr(attrName string) (attr *data.Attribute, exists bool) {

	if task.attrs != nil {
		attr, found := task.attrs[attrName]
		if found {
			return attr, true
		}
	}

	return nil, false
}

// ToLinks returns the predecessor links of the task
func (task *Task) ToLinks() []*Link {
	return task.toLinks
}

// FromLinks returns the successor links of the task
func (task *Task) FromLinks() []*Link {
	return task.fromLinks
}

// InputMapper returns the InputMapper of the task
func (task *Task) InputMapper() *data.Mapper {
	return task.inputMapper
}

// OutputMapper returns the OutputMapper of the task
func (task *Task) OutputMapper() *data.Mapper {
	return task.outputMapper
}

func (task *Task) String() string {
	return fmt.Sprintf("Task[%d]:'%s'", task.id, task.name)
}

////////////////////////////////////////////////////////////////////////////
// Link

// LinkType is an enum for possible Link Types
type LinkType int

const (
	// LtDependency denotes an normal dependency link
	LtDependency LinkType = 0

	// LtExpression denotes a link with an expression
	LtExpression LinkType = 1 //expr language on the model or def?

	// LtLabel denotes 'labelled' link
	LtLabel LinkType = 2
)

// Link is the object that describes the definition of
// a link.
type Link struct {
	id         int
	name       string
	fromTask   *Task
	toTask     *Task
	linkType   LinkType
	value      string //expression or label

	definition *Definition
	parent     *Task
}

func (task *Link) ID() int {
	return task.id
}

// FromTask returns the task the link is coming from
func (link *Link) FromTask() *Task {
	return link.fromTask
}

// ToTask returns the task the link is going to
func (link *Link) ToTask() *Task {
	return link.toTask
}

func (link *Link) String() string {
	return fmt.Sprintf("Link[%d]:'%s' - [from:%d, to:%d]", link.id, link.name, link.fromTask.id, link.toTask.id)
}
