package service

import (
	"context"
	"fmt"
	"github.com/deface90/def-feelings/storage"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"time"
)

type TelegramWorker struct {
	engine storage.Engine
	config storage.Config
	logger *log.Logger

	bot *tgbotapi.BotAPI
}

func NewTelegramWorker(e storage.Engine, c storage.Config, l *log.Logger) (*TelegramWorker, error) {
	worker := TelegramWorker{
		engine: e,
		config: c,
		logger: l,
	}

	var err error
	worker.bot, err = tgbotapi.NewBotAPI(c.TelegramBotToken)
	if err != nil {
		return nil, err
	}

	l.Printf("Authorized on account %s", worker.bot.Self.UserName)

	return &worker, nil
}

func (w *TelegramWorker) Exec() {
	go w.runEventListener()
	ticker := time.NewTicker(time.Minute)
	for {
		<-ticker.C
		w.sendNotifications()
	}
}

func (w *TelegramWorker) sendNotifications() {
	ctx := context.Background()
	subList, err := w.engine.ListUsersSubscriptions(ctx, storage.NotificationTypeTelegram, storage.StatusActive)
	if err != nil {
		log.WithError(err).Error("Failed to get user subscriptions for notification worker")
		return
	}

	for _, sub := range subList {
		if time.Now().Before(sub.LastNotification.Add(time.Duration(sub.Frequency) * time.Minute)) {
			continue
		}

		msg := tgbotapi.NewMessage(sub.ChatID, fmt.Sprintf(`
Hi there! It's time to save your current feelings!
Please, visit %vstatus/create via your browser or just type them right here!
Type feelings comma-sepataed, if you want to send a message for your status, start it with a new line.`, w.config.BaseURL))
		_, err = w.bot.Send(msg)
		if err != nil {
			w.logger.WithError(err).Error("Failed to send user notification")
			sub.Status = storage.StatusInactive
		} else {
			sub.Status = storage.StatusActive
		}

		sub.LastNotification = time.Now()
		err = w.engine.EditUserSubscription(ctx, sub)
		if err != nil {
			log.WithError(err).Error("Failed to update user subscription after notification")
			continue
		}

		log.Printf("Successfull sended notification to user %v", sub.UserID)
	}
}

func (w *TelegramWorker) runEventListener() {
	ctx := context.Background()
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := w.bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			w.logger.Printf("[%v]: %v", update.Message.From.UserName, update.Message.Text)
			switch update.Message.Text {
			default:
				sub, userID, err := w.getUserSubscriptionByUsername(ctx, update.Message.From)
				if err != nil {
					w.sendErrorMessage(update.Message.Chat.ID, err.Error())
					continue
				}

				var msg tgbotapi.MessageConfig
				if sub == nil || sub.Status == storage.StatusInactive {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Please, type '/start' to subscribe for def-feelings notifications")
				} else {
					w.processUserStatus(userID, update.Message.Chat.ID, update.Message.Text)
					continue
				}

				_, err = w.bot.Send(msg)
				if err != nil {
					w.logger.WithError(err).Error("Failed to send reply")
				}
			case "/start":
				sub, userID, err := w.getUserSubscriptionByUsername(ctx, update.Message.From)
				if err != nil {
					w.sendErrorMessage(update.Message.Chat.ID, err.Error())
					continue
				}
				if sub == nil {
					sub = &storage.UserSubscription{
						UserID:           userID,
						Type:             storage.NotificationTypeTelegram,
						LastNotification: time.Now(),
					}
				}
				sub.ChatID = update.Message.Chat.ID
				sub.Status = storage.StatusActive
				_, err = w.engine.CreateOrEditUserSubscription(ctx, sub)
				if err != nil {
					w.logger.WithError(err).Error("Failed to create or update user subscription")
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Thank you, you was subscribed to def-feelings bot notifications!")

				_, err = w.bot.Send(msg)
				if err != nil {
					w.logger.WithError(err).Error("Failed to send reply")
				}
			case "/stop":
				sub, _, err := w.getUserSubscriptionByUsername(ctx, update.Message.From)
				if err != nil {
					w.sendErrorMessage(update.Message.Chat.ID, err.Error())
					continue
				}

				var msg tgbotapi.MessageConfig
				if sub == nil {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Your subscription not found")
				} else {
					sub.Status = storage.StatusInactive
					err = w.engine.EditUserSubscription(ctx, sub)
					if err != nil {
						log.WithError(err).Error("failed to edit user subscription")
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Internal server error, please try again later")
					} else {
						msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Your subscription was disabled")
					}
				}

				_, err = w.bot.Send(msg)
				if err != nil {
					w.logger.WithError(err).Error("Failed to send reply")
				}
			}
		}
	}
}

func (w *TelegramWorker) processUserStatus(userID, chatID int64, msg string) {
	log.Printf("Processing user %v status: %v", userID, msg)
	splitMsg := strings.SplitN(msg, "\n", 2)
	if len(splitMsg) == 0 {
		return
	}

	r := regexp.MustCompile(`[\\p{L}]+`)
	feelings := r.FindAllString(splitMsg[0], -1)

	status := &storage.Status{
		UserID:   userID,
		Feelings: feelings,
		Created:  time.Now(),
	}

	if len(splitMsg) > 1 {
		status.Message = splitMsg[1]
	}

	err := status.Validate()
	if err != nil {
		w.sendErrorMessage(chatID, err.Error())
		return
	}

	_, err = w.engine.CreateStatus(context.Background(), status)
	if err != nil {
		log.WithError(err).Error("Failed to create status from telegram message")
		w.sendErrorMessage(chatID, "")
		return
	}

	replyMsg := tgbotapi.NewMessage(chatID, fmt.Sprintf(`
	Your status successfully saved!
	Check your feelings and statuses on %v`, w.config.BaseURL))
	_, err = w.bot.Send(replyMsg)
	if err != nil {
		log.WithError(err).Error("Failed to send reply to user for his status message")
	}
}

func (w *TelegramWorker) getUserSubscriptionByUsername(ctx context.Context, tgUsername *tgbotapi.User) (*storage.UserSubscription, int64, error) {
	username := tgUsername.UserName
	if username == "" {
		username = fmt.Sprintf("%v", tgUsername.ID)
	}

	user, err := w.engine.GetUserByTgUsername(ctx, username)
	if err != nil {
		w.logger.WithError(err).Errorf("failed to find user whith typed to telegram bot")
		return nil, 0, errors.New("Internal server error, please try again later")
	}
	if user == nil {
		return nil, 0, errors.New("Please, register on def-feelings and specify your telegram username in profile")
	}

	sub, err := w.engine.GetUserSubscription(ctx, user.ID, user.NotificationType)
	if err != nil {
		w.logger.WithError(err).Errorf("failed to find user whith typed to telegram bot")
		return nil, 0, errors.New("Internal server error, please try again later")
	}

	return sub, user.ID, nil
}

func (w *TelegramWorker) sendErrorMessage(chatID int64, errStr string) {
	if errStr == "" {
		errStr = "Internal server error occurred, please try again later"
	}
	msg := tgbotapi.NewMessage(chatID, errStr)

	_, err := w.bot.Send(msg)
	if err != nil {
		w.logger.WithError(err).Error("Failed to send reply")
	}
}
