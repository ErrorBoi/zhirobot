package bot

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) faq(m *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(m.Chat.ID, "Жирофак: http://tiny.cc/qw1vcz")
	msg.ReplyToMessageID = m.MessageID
	b.BotAPI.Send(msg)
}

func (b *Bot) start(m *tgbotapi.Message) {
	b.help(m)
	err := b.DB.CreateUser(m.From.ID)
	if err != nil {
		b.lg.Errorf("Create user error: %w", err)
	}
}

func (b *Bot) help(m *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(m.Chat.ID, helpText)
	msg.ParseMode = tgbotapi.ModeHTML
	b.BotAPI.Send(msg)
}

func (b *Bot) setWeight(m *tgbotapi.Message) {
	args := m.CommandArguments()
	args = strings.TrimSpace(args)
	var msg string

	if len(args) != 0 {
		userWeightStr := strings.Split(args, " ")[0]
		userWeightStr = strings.Replace(userWeightStr, ",", ".", 1)

		if userWeightFloat64, err := strconv.ParseFloat(userWeightStr, 64); err != nil {
			msg = fmt.Sprintf("%s не является корректным числом.", userWeightStr)
		} else {
			if userWeightFloat64 > 0 {
				weightDiff, err := b.DB.SetUserWeight(m.From.ID, userWeightFloat64)
				if err != nil {
					b.lg.Errorf("Set User Weight error: %w", err)
				}
				switch {
				case *weightDiff == userWeightFloat64:
					msg = "Вес записан! (◕‿◕✿)"
				case *weightDiff < 0:
					msg = fmt.Sprintf("Вес записан! (◕‿◕✿)\nС момента прошлого взвешивания сброшено <b>%.1f</b> кг.", math.Abs(*weightDiff))
				case *weightDiff > 0:
					msg = fmt.Sprintf("Вес записан! (◕‿◕✿)\nС момента прошлого взвешивания набрано <b>%.1f</b> кг.", *weightDiff)
				case *weightDiff == 0:
					msg = "Вес записан! (◕‿◕✿)\nС момента прошлого взвешивания вес не изменился."
				}
			} else {
				msg = "Введи положительное число!"
			}
		}
	} else {
		msg = "После команды нужно написать вес, например <pre>/setweight 85.3</pre>"
	}
	message := tgbotapi.NewMessage(m.Chat.ID, msg)
	message.ParseMode = tgbotapi.ModeHTML
	b.BotAPI.Send(message)
}

func (b *Bot) getWeight(m *tgbotapi.Message) {
	stats, err := b.DB.GetUserWeight(m.From.ID)
	if err != nil {
		b.lg.Errorf("Get User Weight error: %w", err)
	}
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
	if err != nil {
		b.lg.Errorf("Get invite link error: %w", err)
	}
	msg := fmt.Sprintf("Инвайт в чат \"Жиросброс\": %s", inviteLink)
	message := tgbotapi.NewMessage(m.Chat.ID, msg)
	b.BotAPI.Send(message)
}

func (b *Bot) turnNotifyOn(m *tgbotapi.Message) {
	err := b.DB.SetNotify(m.From.ID, true)
	if err != nil {
		b.lg.Errorf("Set Notify On error: %w", err)
	}

	msg := fmt.Sprintf("%s, еженедельные уведомления включены. Отключить их можно командой <pre>/off</pre>", m.From.FirstName)
	message := tgbotapi.NewMessage(m.Chat.ID, msg)
	message.ParseMode = tgbotapi.ModeHTML
	b.BotAPI.Send(message)
}

func (b *Bot) turnNotifyOff(m *tgbotapi.Message) {
	err := b.DB.SetNotify(m.From.ID, false)
	if err != nil {
		b.lg.Errorf("Set Notify Off error: %w", err)
	}

	msg := fmt.Sprintf("%s, еженедельные уведомления отключены. Включить их можно командой <pre>/on</pre>", m.From.FirstName)
	message := tgbotapi.NewMessage(m.Chat.ID, msg)
	message.ParseMode = tgbotapi.ModeHTML
	b.BotAPI.Send(message)
}
