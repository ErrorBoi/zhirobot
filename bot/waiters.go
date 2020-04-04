package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
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

func (b *Bot) wakeUp() {
	_, err := http.Get("https://zhirobot.herokuapp.com/")
	if err != nil {
		panic(err)
	}
}

func (b *Bot) testAutoWakeup() {
	ccfg := tgbotapi.ChatConfig{
		ChatID: 128883002,
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

	_, err = b.BotAPI.Send(text)
	if err != nil {
		panic(err)
	}
}