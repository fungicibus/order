package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	"github.com/fungicibus/order/config"
	v1 "github.com/fungicibus/order/internal/api/v1"
	"github.com/fungicibus/order/internal/logger"
	"github.com/fungicibus/order/internal/server"
)

var Tag string
var Commit string

func main() {
	version := getVersion()

	cfg, err := config.GetDefault()
	if err != nil {
		panic(fmt.Errorf("failed to get config: %w", err))
	}
	cfg.App.Version = version

	vmLogs := logger.NewVictoriaLogsWriter(cfg.Log.VictoriaUrl)
	log, err := logger.New(cfg, vmLogs)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get logger")
	}

	cfgContent, _ := json.Marshal(cfg)
	log.Debug().RawJSON("config", cfgContent).Send()

	v1 := v1.New(cfg, log)

	server := server.New(cfg, log, v1.GetHandler())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	go func() {
		err := server.Run(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	<-ctx.Done()
	server.Shutdown()
}

func getVersion() string {
	tag, commit := Tag, Commit

	if Tag == "" {
		tag = "tag"
	}
	if Commit == "" {
		commit = "commit"
	}
	return tag + "-" + commit
}
