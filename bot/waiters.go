package bot

import (
	"fmt"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) weeklyNotification() {
	b.groupWeeklyNotification()
	b.usersWeeklyNotification()
}

func (b *Bot) wakeUp() {
	_, err := http.Get("https://zhirobot.herokuapp.com/")
	if err != nil {
		b.lg.Errorf("Send HTTP request error: %w", err)
	}
}

func (b *Bot) groupWeeklyNotification() {
	text := tgbotapi.NewMessage(b.ChatID, weeklyNotificationMessage)
	message, err := b.BotAPI.Send(text)
	if err != nil {
		b.lg.Errorf("Send message error: %w", err)
	}

	pcmcfg := tgbotapi.PinChatMessageConfig{
		ChatID:              b.ChatID,
		MessageID:           message.MessageID,
		DisableNotification: false,
	}
	_, err = b.BotAPI.PinChatMessage(pcmcfg)
	if err != nil {
		b.lg.Errorf("Pin message error: %w", err)
	}
}

func (b *Bot) usersWeeklyNotification() {
	users, err := b.DB.GetUsers()
	if err != nil {
		b.lg.Errorf("Get users error: %w", err)
	}
	for _, user := range users {
		if !user.Notify {
			continue
		}

		ccfg := tgbotapi.ChatConfig{
			ChatID: int64(user.TgID),
		}
		chat, err := b.BotAPI.GetChat(ccfg)
		if err != nil {
			b.lg.Errorf("Get chat error: %w", err)
		}

		text := tgbotapi.NewMessage(int64(user.TgID), fmt.Sprintf(weeklyUserNotificationMessage, chat.FirstName))
		_, err = b.BotAPI.Send(text)
		if err != nil {
			b.lg.Errorf("Send message error: %w", err)
		}
	}
}
