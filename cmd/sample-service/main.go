package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/andrii2g/go-api-key-gateway/apikey"
	"github.com/andrii2g/go-api-key-gateway/config"
	"github.com/andrii2g/go-api-key-gateway/httpapi"
	"github.com/andrii2g/go-api-key-gateway/storage/sqlstore"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	pepper, err := config.LoadPepper(cfg.PepperBase64, cfg.PepperFile)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := sqlstore.Open(ctx, sqlstore.Options{DSN: cfg.DBDSN})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := sqlstore.NewStore(db)
	usage := apikey.NewAsyncUsageRecorder(store, cfg.UsageQueueSize, cfg.MarkUsedTimeout)
	defer usage.Close(context.Background())

	service, err := apikey.NewService(store, apikey.Options{
		Pepper:            pepper,
		SecretBytes:       cfg.SecretBytes,
		UsageQueueSize:    cfg.UsageQueueSize,
		UsageFlushTimeout: cfg.UsageFlushTimeout,
		MarkUsedTimeout:   cfg.MarkUsedTimeout,
	}, usage)
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: httpapi.NewServer(service, cfg.Environment, cfg.AdminToken),
	}
	log.Printf("sample service listening on %s", cfg.HTTPAddr)
	log.Fatal(server.ListenAndServe())
}
