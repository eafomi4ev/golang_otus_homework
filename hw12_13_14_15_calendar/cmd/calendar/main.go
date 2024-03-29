package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/eafomi4ev/golang_otus_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/eafomi4ev/golang_otus_homework/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/eafomi4ev/golang_otus_homework/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/eafomi4ev/golang_otus_homework/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/eafomi4ev/golang_otus_homework/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/eafomi4ev/golang_otus_homework/hw12_13_14_15_calendar/internal/storage/sql"
)

var appConfigPath string

func init() {
	flag.StringVar(&appConfigPath, "config", "./configs/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	config := NewConfig(appConfigPath)
	configErrors := config.Validate()
	if len(configErrors) != 0 {
		log.Fatal(configErrors)
	}

	logg := logger.New(ctx, config.Logger.Level, config.Logger.Path)

	storage, err := initStorage(ctx, config.Storage)
	if err != nil {
		log.Fatal(fmt.Errorf("error occurred on attempt to crate storage: %w", err))
	}

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(config.Service.Host, config.Service.Port, calendar, logg)
	grpcServer := internalgrpc.NewServer(config.GRPCService.Host, config.GRPCService.Port, calendar, logg)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals)

		<-signals
		signal.Stop(signals)
		cancel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
		grpcServer.Stop()
	}()

	if err := grpcServer.Start(ctx); err != nil {
		logg.Error("failed to start grpc server: " + err.Error())
		os.Exit(1)
	}
	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		os.Exit(1)
	}
}

func initStorage(ctx context.Context, conf StorageConf) (storage app.Storage, err error) {
	switch conf.Type {
	case "inmemory":
		storage = memorystorage.New()
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", conf.Postgres.Host, conf.Postgres.Port, conf.Postgres.User, conf.Postgres.Password, conf.Postgres.DBName, conf.Postgres.SSLMode)
		if storage, err = sqlstorage.New(ctx, dsn); err != nil {
			return nil, fmt.Errorf("cannot init postgres storage: %w", err)
		}
	}

	return storage, nil
}
