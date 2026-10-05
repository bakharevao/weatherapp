package main

import (
	"context"
	"database/sql"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"weatherapp/handler"
	"weatherapp/middlewire"
	"weatherapp/provider"
	"weatherapp/repository"
	"weatherapp/service"
)

const openMeteoBaseURL = "https://api.open-meteo.com/v1/forecast"

func main() {
	logger := logrus.New()

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "root:root@tcp(127.0.0.1:3306)/weatherapp?parseTime=true"
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Fatalf("failed to connect to db: %v", err)
	}

	cityCoordinatesRepository := repository.NewCityCoordinatesRepository(db)
	weatherProvider := provider.NewOpenMeteoProvider(openMeteoBaseURL, cityCoordinatesRepository)
	weatherService := service.NewWeatherService(weatherProvider, cityCoordinatesRepository, 5*time.Second)

	if token := os.Getenv("TELEGRAM_BOT_TOKEN"); token != "" {
		bot, err := tgbotapi.NewBotAPI(token)
		if err != nil {
			logger.Fatalf("failed to create telegram bot: %v", err)
		}

		updateConfig := tgbotapi.NewUpdate(0)
		updateConfig.Timeout = 30
		updates := bot.GetUpdatesChan(updateConfig)

		go func() {
			for update := range updates {
				handler.HandleTelegramUpdate(context.Background(), weatherService, bot, logger, update)
			}
		}()
	} else {
		logger.Warn("TELEGRAM_BOT_TOKEN not set, telegram bot disabled")
	}

	e := echo.New()
	e.Use(middlewire.LoggerMiddleware(logger))

	e.GET("/weather/:city", handler.GetWeatherHandler(weatherService))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logger.Fatal(e.Start(":" + port))
}
