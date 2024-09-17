package filewrite

import (
	"fmt"
	"log"
	"log-aggregation/consumer/types"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/streadway/amqp"
)

func retryWithBackoff(attempt int, maxRetries int) time.Duration {
	if attempt >= maxRetries {
		return 0
	}

	backoff := math.Pow(2, float64(attempt))
	return time.Duration(backoff) * time.Second
}

func ProcessAndRetryWriteLogs(logLevel string, logs []amqp.Delivery, attempt int, maxRetries int) error {
	err := writeBulkToLogFile(logLevel, logs)
	if err == nil {
		return nil
	}

	if attempt < maxRetries {
		waitTime := retryWithBackoff(attempt, maxRetries)
		log.Printf("Retrying in %v (attempt %d/%d)\n", waitTime, attempt+1, maxRetries)
		time.Sleep(waitTime)
		return ProcessAndRetryWriteLogs(logLevel, logs, attempt+1, maxRetries)
	}

	log.Printf("Failed to process logs after %d retries", maxRetries)
	return err
}

func writeBulkToLogFile(logLevel string, messages []amqp.Delivery) error {
	filePath, exists := types.LogFilePaths[logLevel]
	if !exists {
		log.Printf("Unknown log level: %s", logLevel)
		return fmt.Errorf("unknown log level: %s", logLevel)
	}

	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Printf("Directory %s does not exist, creating it...", dir)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %v", dir, err)
			return err
		}
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open log file %s: %v", filePath, err)
		return err
	}
	defer file.Close()

	for _, msg := range messages {
		_, err := file.WriteString(string(msg.Body) + "\n")
		if err != nil {
			log.Printf("Failed to write to log file %s: %v", filePath, err)
			return err
		}
	}

	return nil
}
