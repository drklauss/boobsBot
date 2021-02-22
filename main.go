package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/trace"

	"github.com/drklauss/boobsBot/bot"
	"github.com/drklauss/boobsBot/config"
	"github.com/drklauss/boobsBot/handlers"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
)

func main() {
	defer trace.Stop()
	f, _ := os.Create("out.trace")
	if err := trace.Start(f); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "could not load yml: %v", err)
		os.Exit(1)
	}
	debug := flag.Bool("debug", false, "enable debug")
	flag.Parse()
	if err := config.Load(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "could not load yml: %v", err)
		os.Exit(1)
	}
	initLogger(*debug)

	b, err := bot.New(config.Get(), *debug)
	if err != nil {
		log.Fatalf("could not create bot %v", err)
	}

	b.Handle("/admin", handlers.Empty)
	b.Handle("/debugStart", handlers.Empty)
	b.Handle("/debugStop", handlers.Empty)
	b.Handle("/help", handlers.Help)
	b.Handle("/rate", handlers.Rate)
	b.Handle("/update", handlers.Update)
	b.Handle("/topViewers", handlers.TopViewers)
	b.Handle("/categoriesStat", handlers.CategoriesStat)
	b.Handle("/categories", handlers.Categories)
	b.Handle("/cats", handlers.Categories)
	for _, c := range config.Get().Categories {
		b.Handle("/"+c.Name, handlers.Get)
	}
	b.UseMiddlewares(bot.LogRequest, bot.CheckAdmin)

	b.Run()
}

// initLogger initializes logger.
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
