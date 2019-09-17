package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jasonlvhit/gocron" // Job Scheduling Package
)

// Bot unites botAPI and channels
type Bot struct {
	BotAPI   *tgbotapi.BotAPI
	ChatName string
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

	bot.ChatName = "@grosebros"

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

	// Send "Time to weigh" reminder every Sunday
	gocron.Every(1).Sunday().At("10:00").Do(b.weeklyNotification, b.ChatName)
	gocron.Start()

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.IsCommand() {
			message := update.Message
			b.ExecuteCommand(message)
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		b.BotAPI.Send(msg)
	}
}

// ExecuteCommand distributes commands to go routines
func (b *Bot) ExecuteCommand(m *tgbotapi.Message) {
	command := m.Command()
	log.Printf("Command: %s, Username: %s, ID: %d", command, m.From.UserName, m.From.ID)

	switch command {
	case "faq":
		{
			go b.faq(m)
		}
	case "start":
		{
			go b.start(m)
		}
	default:
		{
			msg := tgbotapi.NewMessage(m.Chat.ID, "Я не знаю такой команды (凸ಠ益ಠ)凸\nНапиши /help и получи справку по командам")
			msg.ReplyToMessageID = m.MessageID
			b.BotAPI.Send(msg)
		}
	}
}
