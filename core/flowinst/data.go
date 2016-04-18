package flowinst

import (
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/flow"
)

func applyInputMapper(pi *Instance, taskData *TaskData) {

	// get the input mapper
	inputMapper := taskData.task.InputMapper()

	if pi.Patch != nil {
		// check if the patch has a overriding mapper
		mapper := pi.Patch.GetInputMapper(taskData.task.ID())
		if mapper != nil {
			inputMapper = mapper
		}
	}

	if inputMapper != nil {
		log.Debug("Applying InputMapper")
		inputMapper.Apply(pi, taskData.InputScope())
	}
}

func applyInputInterceptor(pi *Instance, taskData *TaskData) bool {

	if pi.Interceptor != nil {

		// check if this task as an interceptor
		taskInterceptor := pi.Interceptor.GetTaskInterceptor(taskData.task.ID())

		if taskInterceptor != nil {

			log.Debug("Applying Interceptor")

			if len(taskInterceptor.Inputs) > 0 {
				// override input attributes
				for _, attribute := range taskInterceptor.Inputs {

					log.Debugf("Overriding Attr: %s = %s", attribute.Name, attribute.Value)

					//todo: validation
					taskData.InputScope().SetAttrValue(attribute.Name, attribute.Value)
				}
			}

			// check if we should not evaluate the task
			return !taskInterceptor.Skip
		}
	}

	return true
}

func applyOutputInterceptor(pi *Instance, taskData *TaskData) {

	if pi.Interceptor != nil {

		// check if this task as an interceptor and overrides ouputs
		taskInterceptor := pi.Interceptor.GetTaskInterceptor(taskData.task.ID())
		if taskInterceptor != nil && len(taskInterceptor.Outputs) > 0 {
			// override output attributes
			for _, attribute := range taskInterceptor.Outputs {

				//todo: validation
				taskData.OutputScope().SetAttrValue(attribute.Name, attribute.Value)
			}
		}
	}
}

// applyOutputMapper applies the output mapper, returns flag indicating if
// there was an output mapper
func applyOutputMapper(pi *Instance, taskData *TaskData) bool {

	// get the Output Mapper for the Task if one exists
	outputMapper := taskData.task.OutputMapper()

	if pi.Patch != nil {
		// check if the patch overrides the Output Mapper
		mapper := pi.Patch.GetOutputMapper(taskData.task.ID())
		if mapper != nil {
			outputMapper = mapper
		}
	}

	if outputMapper != nil {
		log.Debug("Applying OutputMapper")
		outputMapper.Apply(taskData.OutputScope(), pi)
		return  true
	}

	return false
}


type FixedTaskScope struct {
	attrs map[string]*data.Attribute
	refAttrs map[string]*data.Attribute
	task *flow.Task
}

func NewFixedTaskScope(refAttrs map[string]*data.Attribute, task *flow.Task) data.Scope {

	scope := &FixedTaskScope{
		refAttrs: refAttrs,
		task: task,
	}

	return scope
}

// GetAttrType implements Scope.GetAttrType
func (s *FixedTaskScope) GetAttrType(attrName string) (attrType string, exists bool) {

	attr, found := s.refAttrs[attrName]

	if found {
		return attr.Type, true
	}

	return "", false
}

// GetAttrValue implements Scope.GetAttrValue
func (s *FixedTaskScope) GetAttrValue(attrName string) (value interface{}, exists bool) {

	if len(s.attrs) > 0 {

		attr, found := s.attrs[attrName]

		if found {
			return attr.Value, true
		}
	}

	if s.task != nil {
		attr, found := s.task.GetAttr(attrName)

		if found {
			return attr.Value, true
		}
	}

	return nil, false
}

// SetAttrValue implements Scope.SetAttrValue
func (s *FixedTaskScope) SetAttrValue(attrName string, value interface{}) {

	if len(s.attrs) == 0 {
		s.attrs = make(map[string]*data.Attribute)
	}

	attr, found := s.attrs[attrName]

	if found {
		attr.Value = value
	} else {
		// look up reference for type
		attr, found = s.refAttrs[attrName]
		if found {
			s.attrs[attrName] = &data.Attribute{Name:attrName, Type:attr.Type, Value:value}
		}
		//todo: else error
	}
}