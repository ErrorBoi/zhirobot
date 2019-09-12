package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Bot unites botAPI and channels
type Bot struct {
	BotAPI *tgbotapi.BotAPI
}

// InitBot inits a bot with given Token
func InitBot(BotToken string) (*Bot, error) {
	var err error
	var bot Bot
	bot.BotAPI, err = tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		return nil, err
	}

	bot.BotAPI.Buffer = 12 * 50

	return &bot, nil
}

// SetDebugMode turns botAPI's debug mode on/off
func (b *Bot) SetDebugMode(DebugMode bool, err error) {
	b.BotAPI.Debug = DebugMode
	if err != nil {
		log.Panic(err)
	}
}

// InitUpdates inits an Updates Channel
func (b *Bot) InitUpdates() {
	ucfg := tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60

	updates, err := b.BotAPI.GetUpdatesChan(ucfg)

	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", b.BotAPI.Self.UserName)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		b.BotAPI.Send(msg)
	}
}
