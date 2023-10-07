package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

// Hostname depicts hostname applicaiton is running on
var Hostname string

// Environment depicts application environment, like development, testing, staging, production and so on
var Environment = "development"

// Address depicts network interface, on which application is listening
var Address = ""

// ListenHTTP shows where HTTP server is binding
var ListenHTTP = ":3000"

// ListenGRPC shows where GRPC server is binding
var ListenGRPC = ":3001"

// Port depicts network port, on which application is listening
var Port = "3000"

// JwtSecret is used to verify JWT tokens
var JwtSecret = "super_secret_for_purser"

// JaegerHost is used to make connection string for JaegerUI for trace collection
var JaegerHost = "127.0.0.1"

// JaegerPort is used to make connection string for JaegerUI for trace collection
var JaegerPort = "6831"

var Domain = "localhost"

var Driver = "memory"
var DatabaseConnectionString string

func IsProduction() bool {
	return Environment == "production"
}

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
}
