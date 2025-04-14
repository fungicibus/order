package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"

	"github.com/fungicibus/order/config"
	v1 "github.com/fungicibus/order/internal/api/v1"
	"github.com/fungicibus/order/internal/logger"
	"github.com/fungicibus/order/internal/server"
)

func main() {
	log := logger.New()

	cfg, err := config.GetDefault()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get config")
	}
	log.SetLevel(cfg.LogLevel)

	prettyJSON, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to marshal config")
	}
	log.Debug().Msgf("Config: /n%s", string(prettyJSON))

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
