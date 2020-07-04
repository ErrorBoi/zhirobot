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
	_, err := b.BotAPI.Send(msg)
	fmt.Println(err)
}

func (b *Bot) setWeight(m *tgbotapi.Message) {
	args := m.CommandArguments()
	b.parseAndSetWeight(m, args)
}

func (b *Bot) setHeight(m *tgbotapi.Message) {
	args := m.CommandArguments()
	b.parseAndSetHeight(m, args)
}

func (b *Bot) getWeight(m *tgbotapi.Message) {
	stats, err := b.DB.GetUserWeight(m.From.ID)
	if err != nil {
		b.lg.Errorf("Get User Weight error: %w", err)
	}
	var imt float64

	height, err := b.DB.GetUserHeight(m.From.ID)
	if err != nil {
		b.lg.Errorf("Get User Height error: %w", err)
		imt = 0
	} else {
		imt, err = b.DB.GetUserIMT(m.From.ID)
		if err != nil {
			b.lg.Errorf("Get User IMT error: %w", err)
			imt = 0
		}
	}

	msg := fmt.Sprintf(`<pre>
%s:
Рост: %d см
ИМТ: %6.1f
|   Вес     |     Дата      |
|-----------|:-------------:|`, m.From.FirstName, height, imt)
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

	msg := fmt.Sprintf("%s, еженедельные уведомления включены. Отключить их можно командой /off", m.From.FirstName)
	message := tgbotapi.NewMessage(int64(m.From.ID), msg)
	b.BotAPI.Send(message)
}

func (b *Bot) turnNotifyOff(m *tgbotapi.Message) {
	err := b.DB.SetNotify(m.From.ID, false)
	if err != nil {
		b.lg.Errorf("Set Notify Off error: %w", err)
	}

	msg := fmt.Sprintf("%s, еженедельные уведомления отключены. Включить их можно командой /on", m.From.FirstName)
	message := tgbotapi.NewMessage(int64(m.From.ID), msg)
	b.BotAPI.Send(message)
}

func (b *Bot) parseAndSetWeight(m *tgbotapi.Message, weight string) {
	weight = strings.TrimSpace(weight)
	var msg string

	if len(weight) != 0 {
		userWeightStr := strings.Split(weight, " ")[0]
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

func (b *Bot) parseAndSetHeight(m *tgbotapi.Message, height string) {
	height = strings.TrimSpace(height)
	var msg string

	if len(height) != 0 {
		userHeightStr := strings.Split(height, " ")[0]
		userHeightStr = strings.Replace(userHeightStr, ",", ".", 1)

		if userHeightInteger, err := strconv.Atoi(userHeightStr); err != nil {
			msg = fmt.Sprintf("%s не является корректным числом.", userHeightStr)
		} else {
			if userHeightInteger > 0 {
				err := b.DB.SetUserHeight(m.From.ID, userHeightInteger)
				if err != nil {
					b.lg.Errorf("Set User Height error: %w", err)
				}
				msg = "Рост записан! (◕‿◕✿)"
			} else {
				msg = "Введи положительное число!"
			}
		}
	} else {
		msg = "После команды нужно написать рост в сантиметрах, например <pre>/setheight 175</pre> или <pre>/sh 175</pre>"
	}
	message := tgbotapi.NewMessage(m.Chat.ID, msg)
	message.ParseMode = tgbotapi.ModeHTML
	b.BotAPI.Send(message)
}
