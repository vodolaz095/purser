package config

import "time"

// LogOutputType задаёт куда и как выводятся логи
type LogOutputType string

const (
	LogOutputConsole    LogOutputType = "console"
	LogOutputStdOutJSON LogOutputType = "stdout_json"
	LogOutputJournald   LogOutputType = "journald"
)

const PruneOldSecretsInterval = 30 * time.Second
