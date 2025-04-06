package main

import (
	"context"
	justMealHttpApi "just-meal/api/http"
	"just-meal/internal/repositories"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	databaseCfg := repositories.DBConfig{
		Type: repositories.Postgres,
		Postgres: repositories.PgConfig{
			URL: "postgres://user:pass@localhost/db",
		},
	}

	repository, err := repositories.NewDishRepository(context.Background(), databaseCfg)
	if err != nil {
		panic(err)
	}

	httpServer := justMealHttpApi.NewWebApiServer(repository)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: httpServer.Router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}
}
