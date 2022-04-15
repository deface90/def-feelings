package main

import (
	"context"
	"github.com/deface90/def-feelings/rest"
	"github.com/deface90/def-feelings/storage"
	"github.com/deface90/def-feelings/storage/adapter"
	"github.com/deface90/go-logger/filename"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jinzhu/configor"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	var cfg storage.Config
	err := configor.Load(&cfg, "config.json")
	if err != nil {
		log.WithError(err).Warnf("Failed to read config.json, using default config values")
	}

	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})

	filenameHook := filename.NewHook()
	filenameHook.Field = "line"
	logger.AddHook(filenameHook)

	engine, err := makeEngine(cfg)
	if err != nil {
		log.Fatalf("Failed to make engine")
	}

	/*tgWorker, err := service.NewTelegramWorker(engine, cfg, logger)
	if err != nil {
		log.WithError(err).Error("Failed to init telegram worker")
	} else {
		go tgWorker.Exec()
	}*/

	restService := rest.NewRestService(engine, cfg, logger)
	restService.Run()
}

func makeEngine(cfg storage.Config) (engine storage.Engine, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if cfg.Storage.Type == "postgres" {
		m, err := migrate.New("file://migrations", cfg.Storage.Postgres.DSN)
		if err != nil {
			log.WithError(err).Fatalf("Failed to init migrations")
		}
		err = m.Up()
		if err != nil && err != migrate.ErrNoChange {
			log.WithError(err).Errorf("Failed to apply migrations")
		}

		var conn *pgxpool.Pool
		conn, err = pgxpool.Connect(ctx, cfg.Storage.Postgres.DSN)
		if err != nil {
			return nil, err
		}

		return adapter.NewPostgres(conn, cfg.Timezone)
	}

	return nil, errors.New("Unknown engine type")
}
