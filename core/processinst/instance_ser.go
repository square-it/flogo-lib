package processinst

import (
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/util"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////
// Process Instance Serialization

type serInstance struct {
	ID          string            `json:"id"`
	Status      Status            `json:"status"`
	State       int               `json:"state"`
	ProcessURI  string            `json:"processUri"`
	Attrs       []*data.Attribute `json:"attrs"`
	WorkQueue   []*WorkItem       `json:"workQueue"`
	RootTaskEnv *TaskEnv          `json:"rootTaskEnv"`
}

// MarshalJSON overrides the default MarshalJSON for ProcessInstance
func (pi *Instance) MarshalJSON() ([]byte, error) {

	queue := make([]*WorkItem, pi.WorkItemQueue.List.Len())

	for i, e := 0, pi.WorkItemQueue.List.Front(); e != nil; i, e = i+1, e.Next() {
		queue[i], _ = e.Value.(*WorkItem)
	}

	attrs := make([]*data.Attribute, 0, len(pi.Attrs))

	for _, value := range pi.Attrs {
		attrs = append(attrs, value)
	}

	return json.Marshal(&serInstance{
		ID:          pi.id,
		Status:      pi.status,
		State:       pi.state,
		Attrs:       attrs,
		ProcessURI:  pi.ProcessURI,
		WorkQueue:   queue,
		RootTaskEnv: pi.RootTaskEnv,
	})
}

// UnmarshalJSON overrides the default UnmarshalJSON for ProcessInstance
func (pi *Instance) UnmarshalJSON(d []byte) error {

	//if pi.processProvider == nil {
	//	panic("process.Provider not specified, required for unmarshalling")
	//}

	ser := &serInstance{}
	if err := json.Unmarshal(d, ser); err != nil {
		return err
	}

	pi.id = ser.ID
	pi.status = ser.Status
	pi.state = ser.State

	pi.ProcessURI = ser.ProcessURI
	//pi.Process = pi.processProvider.GetProcess(pi.ProcessURI)
	//pi.ProcessModel = processmodel.Get(pi.Process.ModelID())

	pi.Attrs = make(map[string]*data.Attribute)

	for _, value := range ser.Attrs {
		pi.Attrs[value.Name] = value
	}

	pi.RootTaskEnv = ser.RootTaskEnv
	//pi.RootTaskEnv.init(pi)

	pi.WorkItemQueue = util.NewQueue()

	for _, workItem := range ser.WorkQueue {
		workItem.TaskData = pi.RootTaskEnv.TaskDatas[workItem.TaskID]
		pi.WorkItemQueue.Push(workItem)
	}

	pi.ChangeTracker = NewInstanceChangeTracker()

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// Task Env Serialization

// MarshalJSON overrides the default MarshalJSON for TaskEnv
func (te *TaskEnv) MarshalJSON() ([]byte, error) {

	t := make([]*TaskData, 0, len(te.TaskDatas))

	for _, value := range te.TaskDatas {
		t = append(t, value)
	}

	l := make([]*LinkData, 0, len(te.LinkDatas))

	for _, value := range te.LinkDatas {
		l = append(l, value)
	}

	return json.Marshal(&struct {
		ID        int         `json:"id"`
		TaskID    int         `json:"taskId"`
		TaskDatas []*TaskData `json:"taskDatas"`
		LinkDatas []*LinkData `json:"linkDatas"`
	}{
		ID:        te.ID,
		TaskID:    te.taskID,
		TaskDatas: t,
		LinkDatas: l,
	})
}

// UnmarshalJSON overrides the default UnmarshalJSON for TaskEnv
func (te *TaskEnv) UnmarshalJSON(data []byte) error {

	ser := &struct {
		ID        int         `json:"id"`
		TaskID    int         `json:"taskId"`
		TaskDatas []*TaskData `json:"taskDatas"`
		LinkDatas []*LinkData `json:"linkDatas"`
	}{}

	if err := json.Unmarshal(data, ser); err != nil {
		return err
	}

	te.ID = ser.ID
	te.taskID = ser.TaskID
	te.TaskDatas = make(map[int]*TaskData)
	te.LinkDatas = make(map[int]*LinkData)

	for _, value := range ser.TaskDatas {
		te.TaskDatas[value.taskID] = value
	}

	for _, value := range ser.LinkDatas {
		te.LinkDatas[value.linkID] = value
	}

	return nil
}

// MarshalJSON overrides the default MarshalJSON for TaskData
func (td *TaskData) MarshalJSON() ([]byte, error) {

	attrs := make([]*data.Attribute, 0, len(td.attrs))

	for _, value := range td.attrs {
		attrs = append(attrs, value)
	}

	return json.Marshal(&struct {
		TaskID int               `json:"taskId"`
		State  int               `json:"state"`
		Attrs  []*data.Attribute `json:"attrs"`
	}{
		TaskID: td.task.ID(),
		State:  td.state,
		Attrs:  attrs,
	})
}

// UnmarshalJSON overrides the default UnmarshalJSON for TaskData
func (td *TaskData) UnmarshalJSON(d []byte) error {
	ser := &struct {
		TaskID int               `json:"taskId"`
		State  int               `json:"state"`
		Attrs  []*data.Attribute `json:"attrs"`
	}{}

	if err := json.Unmarshal(d, ser); err != nil {
		return err
	}

	td.state = ser.State
	td.taskID = ser.TaskID

	if ser.Attrs != nil {
		td.attrs = make(map[string]*data.Attribute)

		for _, value := range ser.Attrs {
			td.attrs[value.Name] = value
		}
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
// Process Instance Changes Serialization

// MarshalJSON overrides the default MarshalJSON for InstanceChangeTracker
func (ict *InstanceChangeTracker) MarshalJSON() ([]byte, error) {

	var wqc []*WorkItemQueueChange

	if ict.wiqChanges != nil {
		wqc = make([]*WorkItemQueueChange, 0, len(ict.wiqChanges))

		for _, value := range ict.wiqChanges {
			wqc = append(wqc, value)
		}

	} else {
		wqc = nil
	}

	var tdc []*TaskDataChange

	if ict.tdChanges != nil {
		tdc = make([]*TaskDataChange, 0, len(ict.tdChanges))

		for _, value := range ict.tdChanges {
			tdc = append(tdc, value)
		}
	} else {
		tdc = nil
	}

	return json.Marshal(&struct {
		Status    Status                 `json:"status"`
		State     int                    `json:"state"`
		WqChanges []*WorkItemQueueChange `json:"wqChanges"`
		TdChanges []*TaskDataChange      `json:"tdChanges"`
	}{
		Status:    ict.instChange.Status,
		State:     ict.instChange.State,
		WqChanges: wqc,
		TdChanges: tdc,
	})
}
