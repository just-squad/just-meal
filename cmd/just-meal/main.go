package main

import (
	"context"
	justMealHttpApi "just-meal/api/http"
	"just-meal/internal/app"
	"just-meal/internal/repositories"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func main() {
	setupAppConfiguration()
	var databaseConfig repositories.DBConfig
	err := viper.Sub("dbConfig").Unmarshal(&databaseConfig)
	if err != nil {
		log.Fatalf("Error while get database config from store: %s", err)
	}
	var appConfig app.AppConfig
	err = viper.Sub("appConfig").Unmarshal(&appConfig)
	if err != nil {
		log.Fatalf("Error while get app config from store: %s", err)
	}
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

func setupAppConfiguration() {
	viper.SetConfigName("appsettings")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("JUSTMEAL")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file: %s", err)
	}
}
