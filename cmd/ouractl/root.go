package ouractl

import (
	"context"
	slog "log"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/fang"
	"github.com/hagelstam/ouractl/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ouractl",
	Short: "Oura Ring CLI",
}

func Execute() error {
	return fang.Execute(context.Background(), rootCmd, fang.WithVersion(version.Version))
}

func init() {
	var loggerFile *os.File
	_, debug := os.LookupEnv("DEBUG")

	if debug {
		var fileErr error
		newLogFile, fileErr := os.OpenFile("debug.log",
			os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
		if fileErr == nil {
			log.SetColorProfile(colorprofile.TrueColor)
			log.SetOutput(newLogFile)
			log.SetTimeFormat("15:04:05.000")
			log.SetReportCaller(true)
			setDebugLogLevel()
			log.Debug("Logging to debug.log")
		} else {
			loggerFile, _ = tea.LogToFile("debug.log", "debug")
			slog.Print("Failed setting up logging", fileErr)
		}
	} else {
		log.SetOutput(os.Stderr)
		log.SetLevel(log.FatalLevel)
	}

	if loggerFile != nil {
		defer loggerFile.Close()
	}
}

func setDebugLogLevel() {
	switch os.Getenv("LOG_LEVEL") {
	case "debug", "":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	}

	log.Debug("log level set", "level", log.GetLevel())
}
