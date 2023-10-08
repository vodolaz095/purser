package config

// LogOutput задаёт куда и как выводятся логи
type LogOutputType string

const (
	LogOutputConsole    LogOutputType = "console"
	LogOutputStdOutJSON LogOutputType = "stdout_json"
	LogOutputJournald   LogOutputType = "journald"
)
