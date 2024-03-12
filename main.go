package main

import (
	"notification-service/app"
	"notification-service/auth"
	"notification-service/config"
	"notification-service/database"
	"notification-service/handlers"
	"notification-service/notification"
	"notification-service/sender"
	"os"
	"os/signal"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sLog := logger.Sugar()

	cfg := config.Config{}
	// cfg.Default()

	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		sLog.Fatalf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		sLog.Fatalf("unmarshal: %v", err)
	}

	db, err := database.NewDB(cfg.DB.URL)
	if err != nil {
		sLog.Fatal(err)
	}
	app := app.NewApp(db, sLog, cfg)
	for nType, sType := range cfg.Senders {
		sender, err := sender.NewSender(sender.SenderType(sType), cfg)
		if err != nil {
			sLog.Fatalf("error creating sender: %v", err)
			continue
		}
		app.AddSender(notification.NotificationTypeFromString(nType), sender)
	}
	sLog.Infof("configured senders: %v", cfg.Senders)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			app.Shutdown()
			os.Exit(1)
		}
	}()
	auth := auth.NewHTTPAPIKeyAuth(cfg.Auth.APIKey)
	handler := handlers.NewHTTP(app, cfg.HTTP.Port, sLog, auth, cfg.HTTP.RatelimitMaxRequests, cfg.HTTP.RatelimitTimeout)
	handler.InitRoutes()
	handler.Run()
}
