package bot

import (
	"log"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/ErrorBoi/zhirobot/db"
)

// Bot unites botAPI and channels
type Bot struct {
	BotAPI *tgbotapi.BotAPI
	DB     *db.DB
	ChatID int64
}

// InitBot inits a bot with given Token
func InitBot(BotToken string, DB *db.DB) (*Bot, error) {
	var err error
	var bot Bot
	bot.BotAPI, err = tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		return nil, err
	}

	bot.BotAPI.Buffer = 12 * 50

	bot.ChatID = ZhirosbrosChatID

	bot.DB = DB

	return &bot, nil
}

// SetDebugMode turns botAPI's debug mode on/off
func (b *Bot) SetDebugMode(DebugMode bool, err error) {
	b.BotAPI.Debug = DebugMode
	if err != nil {
		panic(err)
	}
}

// InitUpdates inits an Updates Channel
func (b *Bot) InitUpdates(BotToken string) {
	ucfg := tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60

	// updates, err := b.BotAPI.GetUpdatesChan(ucfg)
	// h.PanicIfErr(err)
	updates := b.BotAPI.ListenForWebhook("/" + BotToken)
	log.Printf("Authorized on account %s", b.BotAPI.Self.UserName)

	go b.RunScheduler()

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.IsCommand() {
			b.ExecuteCommand(update.Message)
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	}
}

// ExecuteCommand distributes commands to go routines
func (b *Bot) ExecuteCommand(m *tgbotapi.Message) {
	command := strings.ToLower(m.Command())

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
	case "setweight", "sw", "s":
		{
			go b.setWeight(m)
		}
	case "getweight", "gw":
		{
			go b.getWeight(m)
		}
	case "invite":
		{
			go b.getInviteLink(m)
		}
	default:
		{
			if m.Chat.IsPrivate() {
				msg := tgbotapi.NewMessage(m.Chat.ID, "Я не знаю такой команды (凸ಠ益ಠ)凸\nНапиши /help и получи справку по командам")
				msg.ReplyToMessageID = m.MessageID
				b.BotAPI.Send(msg)
			}
		}
	}
}

func (b *Bot) RunScheduler() {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic(err)
	}

	scheduler := gocron.NewScheduler(loc)

	// Send "Time to weigh" reminder every Sunday
	scheduler.Every(1).Sunday().At("11:00").Do(b.weeklyNotification)

	// Wake Up a bot before it goes to idling
	scheduler.Every(15).Minute().Do(b.wakeUp)

	<-scheduler.Start()
}
