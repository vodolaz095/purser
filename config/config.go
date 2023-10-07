package config

import (
	"log"
	"os"
)

// Hostname depicts hostname applicaiton is running on
var Hostname string

// Environment depicts application environment, like development, testing, staging, production and so on
var Environment = "development"

// Address depicts network interface, on which application is listening
var Address = "0.0.0.0"

// Port depicts network port, on which application is listening
var Port = "3000"

// JaegerHost is used to make connection string for JaegerUI for trace collection
var JaegerHost = "127.0.0.1"

// JaegerPort is used to make connection string for JaegerUI for trace collection
var JaegerPort = "6831"

var Driver = "memory"
var DatabaseConnectionString string

func init() {
	var err error
	loadFromEnvironment(&Hostname, "HOSTNAME")
	if Hostname == "" {
		Hostname, err = os.Hostname()
		if err != nil {
			log.Fatalf("error finding hostname: %s", err)
		}
	}
	loadFromEnvironment(&Port, "PORT")
	loadFromEnvironment(&Address, "ADDR")

	loadFromEnvironment(&JaegerHost, "JAEGER_HOST")
	loadFromEnvironment(&JaegerPort, "JAEGER_PORT")

	loadFromEnvironment(&Driver, "DRIVER")
	loadFromEnvironment(&DatabaseConnectionString, "DB_URL")
}
