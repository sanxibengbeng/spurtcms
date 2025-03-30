package loggger

import (
	"fmt"
	"log"
	"os"
)

var (
	LogFilePath = "graphql/logs/errors.log"

	graphqlLogger *GraphQLLogger
)

const (
	InfoLevel = iota
	WarningLevel
	ErrorLevel
)

type GraphQLLogger struct {
	Level       int
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
}

func init() {

	_, err := os.Stat(LogFilePath)

	var LogFile *os.File

	if os.IsNotExist(err) {

		makErr := os.MkdirAll("graphql/log/", 0777)

		if makErr != nil {
			fmt.Printf("Error creating logs directory: %v\n", makErr)
		}

		_, err = os.Create(LogFilePath)

		if err != nil {
			fmt.Printf("Failed to create log file: %v\n", err)
		}

		fmt.Println("Created new log file")
	}

	LogFile, err = os.OpenFile(LogFilePath, os.O_APPEND|os.O_WRONLY, 0666)

	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
	}

	graphqlLogger = &GraphQLLogger{
		Level:       InfoLevel,
		InfoLogger:  log.New(LogFile, "INFO ", log.LstdFlags|log.Lshortfile),
		WarnLogger:  log.New(LogFile, "WARN ", log.LstdFlags|log.Lshortfile),
		ErrorLogger: log.New(LogFile, "ERROR ", log.LstdFlags|log.Lshortfile),
	}

}

func ErrorLog() *log.Logger {
	return graphqlLogger.ErrorLogger
}

func WarnLog() *log.Logger {
	return graphqlLogger.ErrorLogger
}
