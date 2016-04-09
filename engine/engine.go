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

	engine.instManager = processinst.NewManager(env.ProcessProvider(), &engine)

	return &engine
}

// Start will start the engine, by starting all of its triggers and runner
func (e *Engine) Start() {

	log.Info("Starting Engine...")

	triggers := trigger.Triggers()

	// initialize triggers
	for _, trigger := range triggers {

		triggerConfig := e.env.engineConfig.Triggers[trigger.Metadata().ID]
		trigger.Init(nil, triggerConfig)
	}

	//start runner
	e.runner.Start()

	// start triggers
	for _, trigger := range triggers {

		log.Debugf("Starting trigger: %s", trigger.Metadata().ID)
		trigger.Start()
	}

	tester := e.env.engineTester
	testerConfig := e.env.engineConfig.TesterConfig

	if tester != nil && testerConfig.Enabled {
		log.Info("Starting Engine Tester...")
		tester.Init(e.instManager, e.runner, testerConfig.Settings)
		tester.Start()
		log.Info("Started Engine Tester...")
	}

	log.Info("Started Engine")

}

// Stop will stop the engine, by stopping all of its triggers and runner
func (e *Engine) Stop() {

	log.Info("Stopping Engine...")

	triggers := trigger.Triggers()

	// stop triggers
	for _, trigger := range triggers {
		log.Debugf("Stopping trigger: %s", trigger.Metadata().ID)
		trigger.Stop()
	}

	tester := e.env.engineTester
	testerConfig := e.env.engineConfig.TesterConfig

	if tester != nil && testerConfig.Enabled {
		log.Info("Stopping Engine Tester...")
		tester.Stop()
		log.Info("Stopped Engine Tester...")
	}

	// stop runner
	e.runner.Stop()

	log.Info("Stopped Engine")
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
