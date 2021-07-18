package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iKOPKACtraxa/otus-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/iKOPKACtraxa/otus-hw/hw12_13_14_15_calendar/internal/logger"
	internalHTTP "github.com/iKOPKACtraxa/otus-hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/iKOPKACtraxa/otus-hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/iKOPKACtraxa/otus-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/iKOPKACtraxa/otus-hw/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFilePath string

const (
	memory = "memory"
	sql    = "sql"
)

func init() {
	flag.StringVar(&configFilePath, "config", "../../configs/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig()
	logg := logger.New(config.Logger.File, config.Logger.Level)
	var storage storage.Storage
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	switch config.Storage.Type {
	case memory:
		storage = memorystorage.New()
	case sql:
		storage, err = sqlstorage.New(ctx, config.Storage.ConnStr, logg)
	default:
		logg.Error(fmt.Sprintf("wrong config.Storage.Type = \"%v\", must be \"%v\" or \"%v\"", config.Storage.Type, memory, sql))
	}
	if err != nil {
		logg.Error("at sql storage creating has got an error: " + err.Error())
	}
	calendar := app.New(logg, storage)

	server := internalHTTP.NewServer(calendar, config.HTTPServer.HostPort)
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGHUP)

		select {
		case <-ctx.Done():
			return
		case <-signals:
		}

		signal.Stop(signals)
		cancel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")
	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
