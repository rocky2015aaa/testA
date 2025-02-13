package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var (
	Date    = ""
	Version = "dev"
	Build   = "dev"
)

const (
	EnvHashOcrSvrPort     = "HASHOCR_SVR_PORT"
	EnvHashOcrSvrLogLevel = "HASHOCR_SVR_LOG_LEVEL"
	EnvHashOcrSvrGinMode  = "HASHOCR_SVR_GIN_MODE"

	EnvFile = ".env"
)

func init() {
	fmt.Printf("Build Date: %s\nBuild Version: %s\nBuild: %s\n\n", Date, Version, Build)
	err := godotenv.Load(EnvFile)
	if err != nil {
		log.Fatalf("Error loading %s file: %v", EnvFile, err)
	}
	logLevel, err := log.ParseLevel(os.Getenv(EnvHashOcrSvrLogLevel))
	if err != nil {
		logLevel = log.DebugLevel
	}
	log.SetLevel(logLevel)
	log.SetFormatter(&log.JSONFormatter{})
}
