package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) weeklyNotification() {
	ccfg := tgbotapi.ChatConfig{
		ChatID: b.ChatID,
	}
	chat, err := b.BotAPI.GetChat(ccfg)
	if err != nil {
		panic(err)
	}
	text := tgbotapi.NewMessage(chat.ID, `Еженедельная сдача показаний!
	1. Взвесься на голодный желудок
	2. Введи /setweight <вес> в конфе или чате с ботом
	3. ...
	4. Ты лучше всех!`)

	message, err := b.BotAPI.Send(text)
	if err != nil {
		panic(err)
	}

	pcmcfg := tgbotapi.PinChatMessageConfig{
		ChatID:              chat.ID,
		MessageID:           message.MessageID,
		DisableNotification: false,
	}
	_, err = b.BotAPI.PinChatMessage(pcmcfg)
	if err != nil {
		panic(err)
	}
}
