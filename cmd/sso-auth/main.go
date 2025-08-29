package main

import (
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"sso-auth/internal/app"
	"sso-auth/internal/config"
	"syscall"
)

const (
	envLocal = "local"
	envDev = "dev"
	envProd = "prod"
)

func main() {

	flag.Parse()
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting application", slog.String("env", cfg.Env))

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <- stop 

	log.Info("stopping aplication", slog.String("signal", sign.String()))

	application.GRPCSrv.Stop()

	log.Info("application stopped")

}



func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal: 
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev: 
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd: 
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
	
