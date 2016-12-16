// This package contains all the default values for the configuration
package config

import (
	"os"
	"strconv"
)

const (
	LOG_LEVEL_DEFAULT         = "DEBUG"
	LOG_LEVEL_KEY             = "LOG_LEVEL"
	RUNNER_WORKERS_KEY        = "RUNNER_WORKERS"
	RUNNER_WORKERS_DEFAULT    = 5
	RUNNER_QUEUE_SIZE_KEY     = "RUNNER_QUEUE_SIZE"
	RUNNER_QUEUE_SIZE_DEFAULT = 50
)

func GetRunnerWorkers() int {
	numWorkers := RUNNER_WORKERS_DEFAULT
	workersEnv := os.Getenv(RUNNER_WORKERS_KEY)
	if len(workersEnv) > 0 {
		i, err := strconv.Atoi(workersEnv)
		if err == nil {
			numWorkers = i
		}
	}
	return numWorkers
}

func GetRunnerQueueSize() int {
	queueSize := RUNNER_QUEUE_SIZE_DEFAULT
	queueSizeEnv := os.Getenv(RUNNER_QUEUE_SIZE_KEY)
	if len(queueSizeEnv) > 0 {
		i, err := strconv.Atoi(queueSizeEnv)
		if err == nil {
			queueSize = i
		}
	}
	return queueSize
}

func GetLogLevel() string {
	logLevelEnv := os.Getenv(LOG_LEVEL_KEY)
	if len(logLevelEnv) > 0 {
		return logLevelEnv
	}
	return LOG_LEVEL_DEFAULT
}
