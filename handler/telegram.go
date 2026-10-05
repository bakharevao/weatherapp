package handler

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"

	"weatherapp/service"
	"weatherapp/types"
)

func HandleTelegramUpdate(ctx context.Context, weatherService service.WeatherService, bot *tgbotapi.BotAPI, logger *logrus.Logger, update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	if update.Message.IsCommand() && update.Message.Command() == "start" {
		handleStart(ctx, weatherService, bot, logger, update.Message.Chat.ID)
		return
	}

	city := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/weather"))
	if city == "" {
		reply(bot, logger, update.Message.Chat.ID, `Send a city name, e.g. "Paris" or "/weather Paris".`)
		return
	}

	stored, err := weatherService.AddWeather(ctx, city)
	if err != nil {
		reply(bot, logger, update.Message.Chat.ID, fmt.Sprintf("Couldn't get weather for %q: %v", city, err))
		return
	}

	logger.WithFields(logrus.Fields{
		"chat_id": update.Message.Chat.ID,
		"city":    city,
		"temp_c":  stored.TempC,
	}).Info("weather lookup succeeded")

	reply(bot, logger, update.Message.Chat.ID, formatWeather(*stored))
}

func handleStart(ctx context.Context, weatherService service.WeatherService, bot *tgbotapi.BotAPI, logger *logrus.Logger, chatID int64) {
	cities, err := weatherService.ListCities(ctx)
	if err != nil || len(cities) == 0 {
		reply(bot, logger, chatID, "Welcome! Send me a city name to get its weather.")
		return
	}

	reply(bot, logger, chatID, fmt.Sprintf(
		"Welcome! Ask me about the weather in one of these cities:\n%s",
		strings.Join(capitalize(cities), "\n"),
	))
}

func capitalize(cities []string) []string {
	out := make([]string, len(cities))
	for i, city := range cities {
		if city == "" {
			continue
		}
		out[i] = strings.ToUpper(city[:1]) + city[1:]
	}
	return out
}

func formatWeather(w types.StoreData) string {
	return fmt.Sprintf("%s: %.1f°C (updated %s)", w.Name, w.TempC, w.LastUpdated)
}

func reply(bot *tgbotapi.BotAPI, logger *logrus.Logger, chatID int64, text string) {
	if _, err := bot.Send(tgbotapi.NewMessage(chatID, text)); err != nil {
		logger.WithError(err).WithField("chat_id", chatID).Error("failed to send telegram message")
		return
	}
	logger.WithField("chat_id", chatID).Info("sent telegram message")
}
