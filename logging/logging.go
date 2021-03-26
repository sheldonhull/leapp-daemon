package logging

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"leapp_daemon/api/util"
	"os"
	"os/user"
)

var logFile *os.File = nil
var context util.Context

func InitializeLogger() {
	err := createLogDir()
	if err != nil {
		logrus.Fatalln("error:", err.Error())
	}

	err = createLogFile()
	if err != nil {
		logrus.Fatalln("error:", err.Error())
	}

	// TODO: export error level, which should depend on the environment
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{ PrettyPrint: true})
	writer := io.MultiWriter(os.Stderr, logFile)
	logrus.SetOutput(writer)
}

func SetContext(ctx *gin.Context) {
	_ = createLogDir()
	logFilePath, _ := getLogFilePath()
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		logFile = nil
		_ = createLogFile()
		writer := io.MultiWriter(os.Stderr, logFile)
		logrus.SetOutput(writer)
	}

	context = util.NewContext(ctx)
}

func Entry() *logrus.Entry {
	_ = createLogDir()
	logFilePath, _ := getLogFilePath()
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		logFile = nil
		_ = createLogFile()
		writer := io.MultiWriter(os.Stderr, logFile)
		logrus.SetOutput(writer)
	}

	return logrus.WithFields(logrus.Fields{
		"requestUri":    context.RequestUri,
		"host":          context.Host,
		"remoteAddress": context.RemoteAddress,
		"method":        context.Method,
		"body":          context.Body,
		"params":        context.Params,
		"header":        context.Header,
	})
}

func Info(args ...interface{}) {
	Entry().Info(args...)
}

func CloseLogFile() {
	_ = logFile.Close()
}

func GetHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}

func createLogFile() error {
	logFilePath, err := getLogFilePath()
	if err != nil {
		return err
	}

	logFile, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Fatalln("error:", err.Error())
	}

	return nil
}

func createLogDir() error {
	dirPath, err := getLogDirPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.Mkdir(dirPath, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func getLogFilePath() (string, error) {
	homeDir, err := GetHomeDir()
	if err != nil {
		return "", err
	}
	logFilePath := fmt.Sprintf("%s/Library/logs/Leapp/daemon/error.log", homeDir)
	return logFilePath, nil
}

func getLogDirPath() (string, error) {
	homeDir, err := GetHomeDir()
	if err != nil {
		return "", err
	}
	dirPath := fmt.Sprintf("%s/Library/logs/Leapp/daemon", homeDir)
	return dirPath, nil
}
