package logging

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"leapp_daemon/controllers/utils"
	"os"
	"os/user"
)

var logFile *os.File = nil

func InitializeLogger() {
	logrus.SetLevel(logrus.ErrorLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{ PrettyPrint: true})

	err := createLogDir()
	if err != nil {
		logrus.Fatalln("error: %s", err.Error())
	}

	err = createLogFile()
	if err != nil {
		logrus.Fatalln("error: %s", err.Error())
	}

	writer := io.MultiWriter(os.Stdout, logFile)
	logrus.SetOutput(writer)
}

func CtxLogger(context *gin.Context) *logrus.Entry {
	_ = createLogDir()
	logFilePath, _ := getLogFilePath()
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		logFile = nil
		_ = createLogFile()
		writer := io.MultiWriter(os.Stdout, logFile)
		logrus.SetOutput(writer)
	}

	contextInfo := utils.NewContext(context)
	return logrus.WithFields(logrus.Fields{
		"requestUri": contextInfo.RequestUri,
		"host": contextInfo.Host,
		"remoteAddress": contextInfo.RemoteAddress,
		"method": contextInfo.Method,
		"body": contextInfo.Body,
		"params": contextInfo.Params,
		"header": contextInfo.Header,
	})
}

func CloseLogFile() {
	_ = logFile.Close()
}

func createLogFile() error {
	logFilePath, err := getLogFilePath()
	if err != nil {
		return err
	}

	logFile, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Fatalln("error: %s", err.Error())
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
	homeDir, err := getHomeDir()
	if err != nil {
		return "", err
	}
	logFilePath := fmt.Sprintf("%s/Library/logs/Leapp/daemon/log.log", homeDir)
	return logFilePath, nil
}

func getLogDirPath() (string, error) {
	homeDir, err := getHomeDir()
	if err != nil {
		return "", err
	}
	dirPath := fmt.Sprintf("%s/Library/logs/Leapp/daemon", homeDir)
	return dirPath, nil
}

func getHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}
