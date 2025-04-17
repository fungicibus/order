package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/fungicibus/order/config"
	v1 "github.com/fungicibus/order/internal/api/v1"
	"github.com/fungicibus/order/internal/logger"
	"github.com/fungicibus/order/internal/server"
	"github.com/fungicibus/order/internal/storage"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

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

	var srv *server.Server
	{
		initCtx, initCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer initCancel()

		postgres, err := storage.New(initCtx, cfg.Postgres)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to init postgres")
		}
		if err := postgres.Ping(initCtx); err != nil {
			log.Fatal().Err(err).Msg("failed to ping postgres")
		}
		if err := postgres.MigrationUp(initCtx, embedMigrations); err != nil {
			log.Fatal().Err(err).Msg("failed to migrate up postgres")
		}

		api := v1.New(cfg, log, postgres)
		srv = server.New(cfg, log, api.GetHandler())
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	go func() {
		err := srv.Run(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	<-ctx.Done()
	srv.Shutdown()
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
