package bot

import (
	"log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jasonlvhit/gocron"
	"go.uber.org/zap"

	"github.com/ErrorBoi/zhirobot/db"
)

// Bot unites botAPI and channels
type Bot struct {
	BotAPI *tgbotapi.BotAPI
	DB     *db.DB
	lg     *zap.SugaredLogger
	ChatID int64
}

// InitBot inits a bot with given Token
func InitBot(BotToken string, DB *db.DB, lg *zap.SugaredLogger) (*Bot, error) {
	var err error
	botAPI, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		return nil, err
	}

	return &Bot{
		BotAPI: botAPI,
		DB:     DB,
		lg:     lg,
		ChatID: ZhirosbrosChatID,
	}, nil
}

// InitUpdates inits an Updates Channel
func (b *Bot) InitUpdates(BotToken string) {
	ucfg := tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60

	updates := b.BotAPI.ListenForWebhook("/" + BotToken)
	log.Printf("Authorized on account %s", b.BotAPI.Self.UserName)

	go b.RunScheduler()

	for update := range updates {
		if update.Message == nil {
			if update.CallbackQuery != nil {
				b.ExecuteCallbackQuery(update.CallbackQuery)
			}
		} else {
			if update.Message.IsCommand() {
				if b.isAllowed(update.Message.From.ID) {
					b.ExecuteCommand(update.Message)
				} else {
					text := "Вы не можете пользоваться ботом, возможно стоит перестать быть 15летней токсичкой"
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
					msg.ReplyToMessageID = update.Message.MessageID
					b.BotAPI.Send(msg)
				}
			} else {
				b.ExecuteText(update.Message)
			}

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		}
	}
}

// ExecuteCommand distributes commands to go routines
func (b *Bot) ExecuteCommand(m *tgbotapi.Message) {
	command := strings.ToLower(m.Command())

	switch command {
	case "faq":
		go b.faq(m)
	case "start":
		go b.start(m)
	case "help":
		go b.help(m)
	case "setweight", "sw":
		go b.setWeight(m)
	case "setheight", "sh":
		go b.setHeight(m)
	case "getweight", "gw":
		msg, last := b.getWeight(m.From.ID, m.From.FirstName, 0)
		message := tgbotapi.NewMessage(m.Chat.ID, msg)
		message.ParseMode = tgbotapi.ModeHTML

		message.ReplyMarkup = b.GetWeightKeyboard(m.From.ID, 0, m.From.FirstName, last)

		b.BotAPI.Send(message)
	case "invite":
		go b.getInviteLink(m)
	case "on":
		go b.turnNotifyOn(m)
	case "off":
		go b.turnNotifyOff(m)
	case "bmi":
		go b.getBMI(m)
	case "repo":
		go b.changeRepoCommand(m)
	default:
		if m.Chat.IsPrivate() {
			msg := tgbotapi.NewMessage(m.Chat.ID, "Я не знаю такой команды (凸ಠ益ಠ)凸\nНапиши /help и получи справку по командам")
			msg.ReplyToMessageID = m.MessageID
			b.BotAPI.Send(msg)
		}
	}
}

// ExecuteCallbackQuery handles callback queries
func (b *Bot) ExecuteCallbackQuery(cq *tgbotapi.CallbackQuery) {
	if strings.HasPrefix(cq.Data, "getWeight") {
		weightInfo := strings.Split(cq.Data, "-")

		tgIDStr := weightInfo[1]
		tgID, err := strconv.Atoi(tgIDStr)
		if err != nil {
			b.lg.Errorf("String to int convertation error: %w", err)
		}

		pageStr := weightInfo[2]
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			b.lg.Errorf("String to int convertation error: %w", err)
		}

		firstName := weightInfo[3]

		msg, last := b.getWeight(tgID, firstName, page)

		message := tgbotapi.NewEditMessageText(cq.Message.Chat.ID, cq.Message.MessageID, msg)
		message.ParseMode = tgbotapi.ModeHTML

		message.ReplyMarkup = b.GetWeightKeyboard(tgID, page, firstName, last)

		b.BotAPI.Send(message)
	}
}

// ExecuteText parses user weight from non-command messages and sends it to database
func (b *Bot) ExecuteText(m *tgbotapi.Message) {
	if m.Chat.IsPrivate() {
		b.parseAndSetWeight(m, m.Text)
	}
}

func (b *Bot) RunScheduler() {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		b.lg.Errorf("Load time location error: %w", err)
	}

	gocron.ChangeLoc(loc)

	// Send "Time to weigh" reminder every Sunday
	gocron.Every(1).Sunday().At("11:00").Do(b.weeklyNotification)

	// Wake Up a bot before it goes to idling
	gocron.Every(15).Minute().Do(b.wakeUp)

	// Start all the pending jobs
	<-gocron.Start()
}
