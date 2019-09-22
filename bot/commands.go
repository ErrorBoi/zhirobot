package bot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ErrorBoi/zhirobot/db"
	h "github.com/ErrorBoi/zhirobot/helpers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) faq(m *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(m.Chat.ID, "Жирофак: http://tiny.cc/qw1vcz")
	msg.ReplyToMessageID = m.MessageID
	b.BotAPI.Send(msg)
}

func (b *Bot) start(m *tgbotapi.Message) {
	b.help(m)
	db.CreateUser(m.From.ID)
}

func (b *Bot) help(m *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(m.Chat.ID, helpText)
	msg.ParseMode = tgbotapi.ModeHTML
	b.BotAPI.Send(msg)
}

func (b *Bot) setWeight(m *tgbotapi.Message) {
	args := m.CommandArguments()
	args = strings.TrimSpace(args)
	msg := "Вес записан! (◕‿◕✿)"

	if len(args) != 0 {
		userWeightStr := strings.Split(args, " ")[0]
		userWeightStr = strings.Replace(userWeightStr, ",", ".", 1)
		var userWeightFloat64 float64
		if h.IsFloat(userWeightStr) {
			var err error
			userWeightFloat64, err = strconv.ParseFloat(userWeightStr, 64)
			h.PanicIfErr(err)
			if userWeightFloat64 > 0 {
				db.SetUserWeight(m.From.ID, userWeightFloat64)
			} else {
				msg = fmt.Sprintln("Введи положительное число!")
			}
		} else {
			msg = fmt.Sprintf("%s не является корректным числом.", userWeightStr)
		}
	} else {
		msg = "После команды нужно написать вес, например <pre>/setweight 85.3</pre>"
	}
	message := tgbotapi.NewMessage(m.Chat.ID, msg)
	message.ParseMode = tgbotapi.ModeHTML
	b.BotAPI.Send(message)
}

func (b *Bot) getWeight(m *tgbotapi.Message) {
	// Add a feature to get stats of all users (who didn't make their stats private?)
	// Add a feature to get stats for a chosen period
	stats := db.GetUserWeight(m.From.ID)
	msg := fmt.Sprintf(`<pre>
%s:
|   Вес     |     Дата      |
|-----------|:-------------:|`, m.From.FirstName)
	for _, stat := range stats {
		msg += fmt.Sprintf("\n|%6.1f     |   %s  |", stat.WeightValue, stat.WeighDate)
	}
	msg += "</pre>"
	message := tgbotapi.NewMessage(m.Chat.ID, msg)
	message.ParseMode = tgbotapi.ModeHTML
	b.BotAPI.Send(message)
}

func (b *Bot) getInviteLink(m *tgbotapi.Message) {
	ccfg := tgbotapi.ChatConfig{
		ChatID: b.ChatID,
	}
	inviteLink, err := b.BotAPI.GetInviteLink(ccfg)
	h.PanicIfErr(err)
	msg := fmt.Sprintf("Инвайт в чат \"Жиросброс\": %s", inviteLink)
	message := tgbotapi.NewMessage(m.Chat.ID, msg)
	b.BotAPI.Send(message)
}
