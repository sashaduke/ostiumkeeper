package main

import (
	"io"
	"log"
	"os"
)

var logger *log.Logger

func init() {
	// Open or create the log file
	logFile, err := os.OpenFile("daemon_log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Create a multi-writer to write to both file and stdout
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// Create a new logger that writes to the multi-writer
	logger = log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}
