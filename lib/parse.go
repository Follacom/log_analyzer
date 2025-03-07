package lib

import (
	"bytes"
	"io"
	"os"
)

func ParseFile(logFile string) bytes.Buffer {
	// Open the log file in read-write mode
	file, err := os.OpenFile(logFile, os.O_RDWR, 0644)
	if err != nil {
		LogError(err)
		return *new(bytes.Buffer)
	}
	// Close the file before finishing the execution of the function
	defer file.Close()

	// Creation of a temporary buffer
	tmp := new(bytes.Buffer)

	// Copy log content to the temp buffer
	if _, err = io.Copy(tmp, file); err != nil {
		LogError(err)
		return *new(bytes.Buffer)
	}

	// Truncate the file
	if err = file.Truncate(0); err != nil {
		LogError(err)
		return *new(bytes.Buffer)
	}

	// Close the file handler
	// it commits the changes
	if err = file.Close(); err != nil {
		LogError(err)
		return *new(bytes.Buffer)
	}

	// Return the buffer as a string
	return *tmp
}
