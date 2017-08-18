package main

import (
	"boobsBot/algorithm"
	"fmt"
	"log"
	"os"
	"boobsBot/config"
)

func main() {

	initLogFile()
	new(algorithm.Dispatcher).Run()
}

// Инициализация лог-файла
func initLogFile() (*os.File, error) {
	file, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return file, fmt.Errorf("Cannot open log file")
	}
	log.SetOutput(file) // 469350825
	log.SetFlags(3)
	log.Println("Bot Parser is starting...")

	return file, nil
}
