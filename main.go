package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/trace"

	"github.com/drklauss/boobsBot/bot"
	"github.com/drklauss/boobsBot/config"
	"github.com/drklauss/boobsBot/handler"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
)

func main() {
	defer trace.Stop()
	f, _ := os.Create("out.trace")
	if err := trace.Start(f); err != nil {
		fmt.Fprintf(os.Stderr, "could not load yml: %v", err)
		os.Exit(1)
	}
	debug := flag.Bool("debug", false, "enable debug")
	flag.Parse()
	if err := config.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "could not load yml: %v", err)
		os.Exit(1)
	}
	initLogger(*debug)
	b, err := bot.New(config.Get())
	if err != nil {
		log.Fatalf("could not create bot %v", err)
	}
	b.Handle("/admin", handler.Empty)
	b.Handle("/debugStart", handler.Empty)
	b.Handle("/debugStop", handler.Empty)
	b.Handle("/help", handler.Help)
	b.Handle("/rate", handler.Rate)
	b.Handle("/update", handler.Update)
	b.Handle("/categories", handler.Categories)
	for _, c := range config.Get().Categories {
		b.Handle(c.Name, handler.Get)
	}
	b.UseMiddlewares(bot.LogRequest, bot.CheckAdmin)
	b.Run()
}

// Logger initializing
func initLogger(debug bool) {
	file, err := os.OpenFile(config.Get().LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		if debug {
			mw := io.MultiWriter(os.Stdout, file)
			log.SetOutput(mw)
		} else {
			log.SetOutput(file)
		}
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	log.SetLevel(log.InfoLevel)
	if debug {
		log.SetLevel(log.DebugLevel)
	}
}
