package engine

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/TIBCOSoftware/flogo-lib/app"
	"github.com/TIBCOSoftware/flogo-lib/config"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/engine/runner"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-lib/util"
)

// Interface for the engine behaviour
type Engine interface {
	Init(directRunner bool) error
	Start() error
	Stop() error
}

// EngineConfig is the type for the Engine Configuration
type EngineConfig struct {
	App            *app.Config
	initialized    bool
	LogLevel       string
	actionRunner   action.Runner
	serviceManager *util.ServiceManager

	triggers map[string]trigger.Trigger
}

// New creates a new Engine
func New(appCfg *app.Config) (Engine, error) {
	// App is required
	if appCfg == nil {
		return nil, errors.New("no App configuration provided")
	}
	// Name is required
	if len(appCfg.Name) == 0 {
		return nil, errors.New("no App name provided")
	}
	// Version is required
	if len(appCfg.Version) == 0 {
		return nil, errors.New("no App version provided")
	}

	//fix up app configuration if it is older
	app.FixUpApp(appCfg)

	logLevel := config.GetLogLevel()

	return &EngineConfig{App: appCfg, serviceManager: util.GetDefaultServiceManager(), LogLevel: logLevel}, nil
}

func (e *EngineConfig) Init(directRunner bool) error {

	if !e.initialized {
		e.initialized = true

		if directRunner {
			e.actionRunner = runner.NewDirect()
		} else {
			runnerConfig := defaultRunnerConfig()
			e.actionRunner = runner.NewPooled(runnerConfig.Pooled)
		}

		propProvider := app.GetPropertyProvider()
		// Initialize the properties
		for id, value := range e.App.Properties {
			propProvider.SetProperty(id, value)
		}

		data.SetPropertyProvider(propProvider)

		actionFactories := action.Factories()
		for _, factory := range actionFactories {
			if initializable, ok := factory.(util.Initializable); ok {

				if err:=initializable.Init(); err != nil {
					return err
				}
			}
		}

		triggers, err := app.CreateTriggers(e.App.Triggers, e.actionRunner)

		if err != nil {
			errorMsg := fmt.Sprintf("Engine: Error Creating trigger instances - %s", err.Error())
			logger.Error(errorMsg)
			panic(errorMsg)
		}

		e.triggers = triggers


		//instanceHelper := app.NewInstanceHelper(e.App, trigger.Factories(), action.Factories())
		//
		//// Create the trigger instances
		//tInstances, err := instanceHelper.CreateTriggers()
		//if err != nil {
		//	errorMsg := fmt.Sprintf("Engine: Error Creating trigger instances - %s", err.Error())
		//	logger.Error(errorMsg)
		//	panic(errorMsg)
		//}
		//
		//// Initialize and register the triggers
		//for key, value := range tInstances {
		//	triggerInterface := value.Interf
		//
		//	//Init
		//	triggerInterface.Init(e.actionRunner)
		//	//Register
		//	trigger.RegisterInstance(key, value)
		//}
		//
		//e.triggers = tInstances
		//
		//// Create the action instances
		//actions, err := instanceHelper.CreateActions()
		//if err != nil {
		//	errorMsg := fmt.Sprintf("Engine: Error Creating action instances - %s", err.Error())
		//	logger.Error(errorMsg)
		//	panic(errorMsg)
		//}
		//
		//// Initialize and register the actions,
		//for key, value := range actions {
		//	action.Register(key, value)
		//	//do we need an init?
		//}
		//
		//e.actions = actions
	}

	return nil
}

//Start initializes and starts the Triggers and initializes the Actions
func (e *EngineConfig) Start() error {
	logger.Info("Engine: Starting...")

	// Todo document RunnerType for engine configuration
	runnerType := config.GetRunnerType()
	e.Init(runnerType == "DIRECT")

	actionRunner := e.actionRunner.(interface{})

	if managedRunner, ok := actionRunner.(util.Managed); ok {
		util.StartManaged("ActionRunner Service", managedRunner)
	}

	logger.Info("Engine: Starting Services...")

	err := e.serviceManager.Start()

	if err != nil {
		logger.Error("Engine: Error Starting Services - " + err.Error())
	} else {
		logger.Info("Engine: Started Services")
	}

	// Start the triggers

	var failed []string

	for key, value := range e.triggers {
		err := util.StartManaged(fmt.Sprintf("Trigger [ '%s' ]", key), value)
		if err != nil {
			logger.Infof("Trigger [%s] failed to start due to error [%s]", key, err.Error())

			//value.Status = trigger.Failed
			//value.Error = err
			logger.Debugf("StackTrace: %s", debug.Stack())
			if config.StopEngineOnError() {
				logger.Debugf("{%s=true}. Stopping engine", config.ENV_STOP_ENGINE_ON_ERROR_KEY)
				logger.Info("Engine: Stopped")
				os.Exit(1)
			}
			failed = append(failed, key)
		} else {
			logger.Infof("Trigger [%s] started", key)
			//value.Status = trigger.Started
		}
	}

	if len(failed) > 0 {
		//remove failed trigger, we have no use for them
		for _, triggerId := range failed {
			delete(e.triggers, triggerId)
		}
	}

	logger.Info("Engine: Started")
	return nil
}

func (e *EngineConfig) Stop() error {
	logger.Info("Engine: Stopping...")

	// Stop Triggers
	for tgrId, tgr := range e.triggers {
		util.StopManaged("Trigger [ "+tgrId+" ]", tgr)
	}

	//tConfigs := e.App.Triggers
	//
	//for _, tConfig := range tConfigs {
	//	// Get instance
	//	tInst := trigger.Instance(tConfig.Id)
	//	if tInst == nil {
	//		//nothing to stop
	//		continue
	//	}
	//	tInterf := tInst.Interf
	//	if tInterf == nil {
	//		//nothing to stop
	//		continue
	//	}
	//	util.StopManaged("Trigger [ "+tConfig.Id+" ]", tInterf)
	//}

	actionRunner := e.actionRunner.(interface{})

	if managedRunner, ok := actionRunner.(util.Managed); ok {
		util.StopManaged("ActionRunner", managedRunner)
	}

	//TODO temporarily add services
	logger.Info("Engine: Stopping Services...")

	err := e.serviceManager.Stop()

	if err != nil {
		logger.Error("Engine: Error Stopping Services - " + err.Error())
	} else {
		logger.Info("Engine: Stopped Services")
	}

	logger.Info("Engine: Stopped")
	return nil
}
