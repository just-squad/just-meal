package main

import (
	"context"
	"errors"
	"fmt"
	justMealHttpApi "just-meal-api/api/http"
	"just-meal-api/internal/app"
	repositories2 "just-meal-api/internal/repositories"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/viper"
)

func main() {
	log.Println("Starting application...")
	setupAppConfiguration()
	appConfig := getAppConfiguration()
	log.Printf("Application name: %s\n", appConfig.App.Name)
	databaseConfig := getDbConfiguration()

	repository, err := repositories2.NewDishRepository(context.Background(), databaseConfig)
	checkErr("Error while create new dish repository", err)

	httpServer := justMealHttpApi.NewWebApiServer(repository)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", appConfig.App.Port),
		Handler: httpServer.Router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
	cfgsPath := "./cfg"
	cfgsType := "yaml"
	err := addConfiguration(cfgsPath, "appsettings", cfgsType, true)
	checkErr("Error while add configuration", err)
	err = addConfiguration(cfgsPath, "appsettings.local", cfgsType, false)
	checkErr("Error while add configuration", err)
}

func addConfiguration(configPath string, configName string, configType string, isRequired bool) error {

	var tempViper *viper.Viper
	appSetLocalFileName := fmt.Sprintf("%s/%s.%s", configPath, configName, configType)
	_, err := os.Stat(appSetLocalFileName)
	if os.IsNotExist(err) && isRequired {
		log.Fatalf("Required configuration file %s does not exists", appSetLocalFileName)
	} else if !os.IsNotExist(err) {
		tempViper = viper.New()
		tempViper.AddConfigPath(configPath)
		tempViper.SetConfigType(configType)
		tempViper.SetConfigName(configName)
		err = tempViper.ReadInConfig()
		checkErr(fmt.Sprintf("Error while add file with name %s and path %s", configName, configPath), err)

		appSetKeys := tempViper.AllKeys()
		for _, key := range appSetKeys {
			keyValue := fmt.Sprintf("%s", tempViper.GetString(key))
			viper.Set(key, keyValue)
		}
	}
	return nil
}

func getAppConfiguration() app.AppConfig {
	var appConfig app.AppConfig
	appConfigSub := viper.Sub("appConfig")
	if appConfigSub == nil {
		log.Fatalf("Not found app configuration in storage")
	}
	err := appConfigSub.Unmarshal(&appConfig)
	checkErr("Error while get app config from store", err)

	return appConfig
}

func getDbConfiguration() repositories2.DBConfig {
	var databaseConfig repositories2.DBConfig
	dbConfigSub := viper.Sub("dbConfig")
	if dbConfigSub == nil {
		log.Fatalf("Not found app configuration in storage")
	}
	err := dbConfigSub.Unmarshal(&databaseConfig)
	checkErr("Error while get database config from store", err)

	return databaseConfig
}

func checkErr(message string, e error) {
	if e != nil {
		log.Fatalf("%s: %s", message, e)
	}
}
