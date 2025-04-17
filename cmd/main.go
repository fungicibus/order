package main

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fungicibus/order/config"
	v1 "github.com/fungicibus/order/internal/api/v1"
	"github.com/fungicibus/order/internal/logger"
	"github.com/fungicibus/order/internal/server"
	"github.com/fungicibus/order/internal/storage"
	"golang.org/x/sync/errgroup"
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

	appCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	g, gCtx := errgroup.WithContext(appCtx)

	g.Go(func() error {
		err := srv.Run(gCtx)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to run server: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Error().Err(err).Msg("service error")
	}

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
