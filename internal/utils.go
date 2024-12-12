package internal

import (
	"fmt"
	"os"
)

// NOTE: This function is for debuging porpusesd
func LogToFile(data string) {
	fileName := "debug.log"
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		os.Exit(1)
	}
	defer file.Close()

	if _, err := file.WriteString(data); err != nil {
		fmt.Println("Error writing to file: ", err)
		os.Exit(1)
	}

}
