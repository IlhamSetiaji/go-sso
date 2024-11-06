package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

func NewLogger(viper *viper.Viper) *log.Logger {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	level := viper.GetString("log.level")
	switch level {
	case "debug":
		logger.SetFlags(log.LstdFlags | log.Lshortfile)
	case "info":
	case "warn", "error", "fatal":
		logger.SetFlags(log.LstdFlags | log.Lshortfile)
	default:
		// Default log level
		logger.SetFlags(log.LstdFlags | log.Lshortfile)
	}

	return logger
}
