package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

// Hostname задает имя сервера, на котором работает приложение
var Hostname string

// Environment задаёт тип окружения приложения - продовое, тестовое и т.д.
var Environment = "development"

// Address задаёт адрес, на котором слушает приложение, по умолчанию - на всех адресах
var Address = ""

// ListenHTTP задаёт адрес, где слушает HTTP сервер
var ListenHTTP = ":3000"

// ListenGRPC задаёт адрес, где слушает GRPC сервер
var ListenGRPC = ":3001"

// Port depicts network port, on which application is listening
var Port = "3000"

// JwtSecret is used to verify JWT tokens
var JwtSecret = "super_secret_for_purser"

// JaegerHost is used to make connection string for JaegerUI for trace collection
var JaegerHost = "127.0.0.1"

// JaegerPort is used to make connection string for JaegerUI for trace collection
var JaegerPort = "6831"

// Domain задаёт домен, на котором работает HTTP и GRPC серверы
var Domain = "localhost"

var Driver = "memory"
var DatabaseConnectionString string

var LogOutput = string(LogOutputConsole)
var LogLevel = "debug"

func IsProduction() bool {
	return Environment == "production"
}

// имхо вообще весь конфиг должен быть таким

// var Addr = ""
// var HttpPort = 3000
// var GrpcPort = 3001
// var EtcdAddr = "http://localhost:2379"
// var EtcdUsername string
// var EtcdPassword string

// всё остальное подгружается из etcd

func init() {
	var err error
	loadFromEnvironment(&Environment, "GO_ENV")

	loadFromEnvironment(&Hostname, "HOSTNAME")
	if Hostname == "" {
		Hostname, err = os.Hostname()
		if err != nil {
			log.Fatalf("error finding hostname: %s", err)
		}
	}
	loadFromEnvironment(&Port, "PORT")
	loadFromEnvironment(&Address, "ADDR")
	loadFromEnvironment(&ListenHTTP, "LISTEN_HTTP")
	port, err := strconv.ParseUint(Port, 10, 64)
	if err != nil {
		log.Fatalf("error parsing port %s as uint64: %s", Port, err)
	}
	if ListenHTTP == "" {
		ListenHTTP = fmt.Sprintf("%s:%v", Address, port)
	}
	loadFromEnvironment(&ListenGRPC, "LISTEN_GRPC")
	if ListenGRPC == "" {
		ListenHTTP = fmt.Sprintf("%s:%v", Address, 1+port)
	}

	loadFromEnvironment(&JaegerHost, "JAEGER_HOST")
	loadFromEnvironment(&JaegerPort, "JAEGER_PORT")

	loadFromEnvironment(&Driver, "DRIVER")
	loadFromEnvironment(&DatabaseConnectionString, "DB_URL")

	loadFromEnvironment(&LogOutput, "LOG_OUTPUT")
	loadFromEnvironment(&LogLevel, "LOG_LEVEL")

}
