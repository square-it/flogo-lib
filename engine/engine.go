package engine

import (
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-lib/core/ext/trigger"
	"github.com/TIBCOSoftware/flogo-lib/core/processinst"
	"github.com/TIBCOSoftware/flogo-lib/engine/runner"
	"github.com/TIBCOSoftware/flogo-lib/util"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("engine")

// Engine creates and executes ProcessInstances.
type Engine struct {
	generator   *util.Generator
	runner      runner.Runner
	env         *Environment
	instManager *processinst.Manager
}

// NewEngine create a new Engine
func NewEngine(env *Environment) *Engine {

	var engine Engine
	engine.generator, _ = util.NewGenerator()
	engine.env = env

	runnerConfig := engine.env.engineConfig.RunnerConfig

	if runnerConfig.Type == "direct" {
		engine.runner = runner.NewDirectRunner(env.stateRecorder, runnerConfig.Direct.MaxStepCount)
	} else {
		engine.runner = runner.NewPooledRunner(runnerConfig.Pooled, env.stateRecorder)
	}

	if log.IsEnabledFor(logging.DEBUG) {
		cfgJSON, _ := json.MarshalIndent(env.engineConfig, "", "  ")
		log.Debugf("Engine Configuration:\n%s\n", string(cfgJSON))
	}

	engine.instManager = processinst.NewManager(env.ProcessProviderService(), &engine)

	return &engine
}

// Start will start the engine, by starting all of its triggers and runner
func (e *Engine) Start() {

	log.Info("Engine: Starting...")

	triggers := trigger.Triggers()

	engineConfig := e.env.engineConfig

	// initialize engine environment
	e.env.Init(e.instManager, e.runner)

	// initialize triggers
	for _, trigger := range triggers {

		triggerConfig := engineConfig.Triggers[trigger.Metadata().ID]
		trigger.Init(nil, triggerConfig)
	}

	// start the process provider service
	processProvider := e.env.ProcessProviderService()
	startManaged("ProcessProvider Service", processProvider)

	// start the state recorder service if enabled
	stateRecorder, enabled := e.env.StateRecorderService()
	if enabled {
		startManaged("StateRecorder Service", stateRecorder)
	}

	startManaged("ProcessRunner Service", e.runner)

	// start triggers
	for _, trigger := range triggers {
		startManaged("Trigger [ "+trigger.Metadata().ID+" ]", trigger)
	}

	// start the engineTester service if enabled
	engineTester, enabled := e.env.EngineTesterService()
	if enabled {
		startManaged("EngineTester Service", engineTester)
	}

	log.Info("Engine: Started")
}

// Stop will stop the engine, by stopping all of its triggers and runner
func (e *Engine) Stop() {

	log.Info("Engine: Stopping...")

	triggers := trigger.Triggers()

	// stop triggers
	for _, trigger := range triggers {
		stopManaged("Trigger [ "+trigger.Metadata().ID+" ]", trigger)
	}

	engineTester, enabled := e.env.EngineTesterService()

	if enabled {
		stopManaged("EngineTester Service", engineTester)
	}

	stopManaged("Process Runner", e.runner)

	stopManaged("ProcessProvider Service", e.env.ProcessProviderService())

	stateRecorder, enabled := e.env.StateRecorderService()

	if enabled {
		stopManaged("StateRecorder Service", stateRecorder)
	}

	log.Info("Engine: Stopped")
}

// NewProcessInstanceID implements processinst.IdGenerator.NewProcessInstanceID
func (e *Engine) NewProcessInstanceID() string {
	return e.generator.NextAsString()
}

// StartProcessInstance implements processinst.IdGenerator.NewProcessInstanceID
func (e *Engine) StartProcessInstance(processURI string, startData map[string]string, replyHandler processinst.ReplyHandler, execOptions *processinst.ExecOptions) string {

	//todo fix for synchronous execution (DirectRunner)

	instance := e.instManager.StartInstance(processURI, startData, replyHandler, execOptions)
	e.runner.RunInstance(instance)

	return instance.ID()
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////

func startManaged(name string, managed util.Managed) {

	log.Debugf("%s: Starting...", name)
	managed.Start()
	log.Debugf("%s: Started", name)
}

func stopManaged(name string, managed util.Managed) {

	log.Debugf("%s: Stopping...", name)

	err := util.StopManaged(managed)

	if err != nil {
		log.Errorf("Error stopping '%s': %s", name, err.Error())
	} else {
		log.Debugf("%s: Stopped", name)
	}
}
