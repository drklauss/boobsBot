package main

import (
	"fmt"
	"log"
	"os"

	"github.com/boobsBot/algorithm"
	"github.com/boobsBot/algorithm/config"
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
	log.SetOutput(file)
	log.SetFlags(3)
	log.Println("-=-=-=-=-=Bot is starting...=-=-=-=-=-")
	return file, nil
}
