package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/trace"

	"github.com/drklauss/boobsBot/bot"
	"github.com/drklauss/boobsBot/config"
	"github.com/drklauss/boobsBot/handler"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/leesper/holmes"
)

func main() {

	defer trace.Stop()
	f, _ := os.Create("out.trace")
	trace.Start(f)
	debug := flag.Bool("debug", false, "enable debug")
	flag.Parse()
	if err := config.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "could not load yml: %v", err)
		os.Exit(1)
	}
	initLogger(debug)
	b, err := bot.New(config.Get())
	if err != nil {
		holmes.Fatalf("could not create bot %v", err)
	}
	b.Handle("/admin", nil)
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
// todo swap logger to logrus
func initLogger(debug *bool) holmes.Logger {
	middlewares := []func(holmes.Logger) holmes.Logger{
		holmes.LogFilePath(config.Get().LogPath),
		holmes.InfoLevel,
	}
	if *debug {
		middlewares = append(middlewares, holmes.DebugLevel)
		middlewares = append(middlewares, holmes.AlsoStdout)
	}

	return holmes.Start(middlewares...)
}
