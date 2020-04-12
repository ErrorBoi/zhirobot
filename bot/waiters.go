package bot

import (
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) weeklyNotification() {
	ccfg := tgbotapi.ChatConfig{
		ChatID: b.ChatID,
	}
	chat, err := b.BotAPI.GetChat(ccfg)
	if err != nil {
		b.lg.Errorf("Get chat error: %w", err)
	}
	text := tgbotapi.NewMessage(chat.ID, weeklyNotificationMessage)

	message, err := b.BotAPI.Send(text)
	if err != nil {
		b.lg.Errorf("Send message error: %w", err)
	}

	pcmcfg := tgbotapi.PinChatMessageConfig{
		ChatID:              chat.ID,
		MessageID:           message.MessageID,
		DisableNotification: false,
	}
	_, err = b.BotAPI.PinChatMessage(pcmcfg)
	if err != nil {
		b.lg.Errorf("Pin message error: %w", err)
	}
}

func (b *Bot) wakeUp() {
	_, err := http.Get("https://zhirobot.herokuapp.com/")
	if err != nil {
		b.lg.Errorf("Send HTTP request error: %w", err)
	}
}
