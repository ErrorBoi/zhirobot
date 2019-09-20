package bot

import (
	"strconv"
	"strings"
	"zhirobot/db"
	h "zhirobot/helpers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) faq(m *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(m.Chat.ID, "Жирофак: http://tiny.cc/qw1vcz")
	msg.ReplyToMessageID = m.MessageID
	b.BotAPI.Send(msg)
}

func (b *Bot) start(m *tgbotapi.Message) {
	db.CreateUser(m.From.ID)
}

func (b *Bot) help(m *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(m.Chat.ID, helpText)
	msg.ParseMode = tgbotapi.ModeHTML
	b.BotAPI.Send(msg)
}

func (b *Bot) setWeight(m *tgbotapi.Message) {
	// add If len(args) != 0
	args := m.CommandArguments()
	args = strings.TrimSpace(args)
	userWeightStr := strings.Split(args, " ")[0]
	userWeightFloat64, err := strconv.ParseFloat(userWeightStr, 64)
	h.PanicIfErr(err)

	db.SetUserWeight(m.From.ID, userWeightFloat64)
}
