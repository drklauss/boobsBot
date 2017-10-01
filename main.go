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

type Flags struct {
	update    int
	statistic string
}

func main() {
	logFile, _ := initLogFile()
	defer logFile.Close()
	var f Flags
	flag.IntVar(&f.update, "u", 0, "Update Links. Example: -u 500 will fetch and save 500 links")
	flag.StringVar(&f.statistic, "s", "", "Statistic. Example: -s top will generate top viewers report")
	flag.Parse()
	if useFlags(f) {
		os.Exit(0)
	}

	// Запуск бота
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

// Запсукает программу с отдельными флагами
func useFlags(f Flags) bool {
	if f.update > 0 {
		provider := new(dataProvider.Provider)
		p := provider.Init(true)
		p.Update(f.update)
		return true
	}
	switch f.statistic {
	case "top":
		p := new(dataProvider.Provider).Init(false)
		b := p.GetTopViewers4Cl()
		fmt.Printf("\n%s\n", b)
		return true
	}

	return false
}
