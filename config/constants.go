package config

import "time"

// LogOutputType задаёт куда и как выводятся логи
type LogOutputType string

const (
	// LogOutputConsole задаёт вывод логов в стандартный вывод в красивом формате с подсветкой
	LogOutputConsole LogOutputType = "console"
	// LogOutputStdOutJSON задаёт вывод логов в стандартный вывод в кодировке JSON
	LogOutputStdOutJSON LogOutputType = "stdout_json"
	// LogOutputJournald задаёт вывод логов в сокет системы журналирования systemd journald
	LogOutputJournald LogOutputType = "journald"
)

// PruneOldSecretsInterval задаёт интервал очистки старых секретов
const PruneOldSecretsInterval = 30 * time.Second
