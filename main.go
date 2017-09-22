package main

import (
	"fmt"
	"log"
	"os"

	"flag"

	"github.com/drklauss/boobsBot/algorithm"
	"github.com/drklauss/boobsBot/algorithm/config"
	"github.com/drklauss/boobsBot/algorithm/dataProvider"
)

func main() {
	logFile, _ := initLogFile()
	defer logFile.Close()

	var up int
	flag.IntVar(&up, "u", 0, "Update Links. Example: -u 500 will fetch and save 500 links")
	flag.Parse()
	if up > 0 {
		provider := new(dataProvider.Provider)
		p := provider.Init()
		p.Update(up)
		os.Exit(0)
	}
	new(algorithm.Dispatcher).Run()

}

// Инициализация лог-файла
func initLogFile() (*os.File, error) {
	file, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return file, fmt.Errorf("Cannot open log file: %v\n", file)
	}
	log.SetOutput(file)
	log.SetFlags(3)
	log.Println("-=-=-=-=Bot is starting=-=-=-=-")
	return file, nil
}
