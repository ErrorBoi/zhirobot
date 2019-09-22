package bot

import (
	h "zhirobot/helpers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) weeklyNotification(chatID int64) {
	ccfg := tgbotapi.ChatConfig{
		ChatID: chatID,
	}
	chat, err := b.BotAPI.GetChat(ccfg)
	h.PanicIfErr(err)
	text := tgbotapi.NewMessage(chat.ID, `Еженедельная сдача показаний!
	1. Взвесься на голодный желудок
	2. Введи /setweight <вес> в конфе или чате с ботом
	3. ...
	4. Ты лучше всех!`)

	message, err := b.BotAPI.Send(text)
	h.PanicIfErr(err)

	pcmcfg := tgbotapi.PinChatMessageConfig{
		ChatID:              chat.ID,
		MessageID:           message.MessageID,
		DisableNotification: false,
	}
	_, err = b.BotAPI.PinChatMessage(pcmcfg)
	h.PanicIfErr(err)
}
