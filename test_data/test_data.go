package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	numLogs     = 1000
	logFilePath = "/home/salman/Desktop/log-aggregation/logfile.log"
)

var (
	logLevels = []string{"INFO", "WARNING", "ERROR", "DEBUG", "FATAL", "PANIC"}
	services  = []string{"auth-service", "payment-service", "notification-service", "user-service"}
	messages  = []string{
		"User logged in successfully",
		"Payment processed",
		"User profile updated",
		"Email sent to user",
		"Error processing payment",
		"Database connection lost",
		"Failed to load user data",
		"User signed out",
		"Session expired",
	}
)

func generateLog() string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logLevel := logLevels[rand.Intn(len(logLevels))]
	service := services[rand.Intn(len(services))]
	message := messages[rand.Intn(len(messages))]
	return fmt.Sprintf("%s [%s] %s: %s\n", timestamp, logLevel, service, message)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	file, err := os.Create(logFilePath)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return
	}
	defer file.Close()

	for i := 0; i < numLogs; i++ {
		logLine := generateLog()
		_, err := file.WriteString(logLine)
		if err != nil {
			fmt.Printf("Failed to write log line: %v\n", err)
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("Dummy logs have been generated successfully!")
}
