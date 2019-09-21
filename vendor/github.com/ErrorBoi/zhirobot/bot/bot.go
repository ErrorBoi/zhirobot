package bot

import (
	"log"
	"net/http"
	"os"
	"strings"

	h "github.com/ErrorBoi/zhirobot/helpers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jasonlvhit/gocron" // Job Scheduling Package
)

// Bot unites botAPI and channels
type Bot struct {
	BotAPI   *tgbotapi.BotAPI
	ChatName string
}

// Handler responds to http request
func Handler(resp http.ResponseWriter, _ *http.Request) {
	resp.Write([]byte("Hi there! I'm Zhirobot!"))
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
	h.PanicIfErr(err)
}

// InitUpdates inits an Updates Channel
func (b *Bot) InitUpdates(BotToken string) {
	// ucfg := tgbotapi.NewUpdate(0)
	// ucfg.Timeout = 60
	ucfg := tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60

	// updates, err := b.BotAPI.GetUpdatesChan(ucfg)
	// h.PanicIfErr(err)
	updates := b.BotAPI.ListenForWebhook("/" + BotToken)
	log.Printf("Authorized on account %s", b.BotAPI.Self.UserName)

	// Send "Time to weigh" reminder every Sunday
	gocron.Every(1).Sunday().At("10:00").Do(b.weeklyNotification, b.ChatName)
	gocron.Start()

	http.HandleFunc("/", Handler)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.IsCommand() {
			message := update.Message
			b.ExecuteCommand(message)
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	}
}

// ExecuteCommand distributes commands to go routines
func (b *Bot) ExecuteCommand(m *tgbotapi.Message) {
	command := strings.ToLower(m.Command())
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
	case "help":
		{
			go b.help(m)
		}
	case "setweight":
		{
			go b.setWeight(m)
		}
	case "getweight":
		{
			go b.getWeight(m)
		}
	default:
		{
			msg := tgbotapi.NewMessage(m.Chat.ID, "Я не знаю такой команды (凸ಠ益ಠ)凸\nНапиши /help и получи справку по командам")
			msg.ReplyToMessageID = m.MessageID
			b.BotAPI.Send(msg)
		}
	}
}
