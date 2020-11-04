package bot

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

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
	b.parseAndSetWeight(m, args)
}

func (b *Bot) setHeight(m *tgbotapi.Message) {
	args := m.CommandArguments()
	b.parseAndSetHeight(m, args)
}

func (b *Bot) getWeight(tgID int, firstName string, page int) (string, bool) {
	stats, last, err := b.DB.GetUserWeight(tgID, page)
	if err != nil {
		b.lg.Errorf("Get User Weight error: %w", err)
	}
	var bmi float64 = 0

	height, err := b.DB.GetUserHeight(tgID)
	if err != nil {
		b.lg.Errorf("Get User Height error: %w", err)
	} else {
		bmi, err = b.DB.GetUserBMI(tgID)
		if err != nil {
			b.lg.Errorf("Get User BMI error: %w", err)
		}
	}

	msg := fmt.Sprintf(`<b>%s</b>:
Рост: %d см
ИМТ: %.1f
<pre>
|   Вес     |     Дата      |
|-----------|:-------------:|`, firstName, height, bmi)
	for _, stat := range stats {
		msg += fmt.Sprintf("\n|%6.1f     |   %s  |", stat.WeightValue, stat.WeighDate)
	}
	msg += "</pre>"

	return msg, *last
}

func (b *Bot) getBMI(m *tgbotapi.Message) {
	BMI, err := b.DB.GetUserBMI(m.From.ID)
	if err != nil {
		b.lg.Errorf("Get User BMI error: %w", err)
	}

	var text string
	switch {
	case BMI <= 16:
		text = "<b>Ярко выраженный дефицит массы тела.</b> Тебе нужно не сбрасывать, а набирать."
	case BMI > 16 && BMI <= 18.5:
		text = "<b>Дефицит массы тела.</b> Возможно стоит задуматься о наборе веса."
	case BMI > 18.5 && BMI <= 25:
		text = "<b>Нормальная масса тела.</b> Meh, fucking normie"
	case BMI > 25 && BMI <= 30:
		text = "<b>Предожирение.</b> Ситуация не критичная, но возможно стоит задуматься о сбросе веса."
	case BMI > 30 && BMI <= 35:
		text = "<b>Ожирение первой степени.</b> Добро пожаловать в жиросброс!"
	case BMI > 35 && BMI <= 40:
		text = "<b>Ожирение второй степени.</b> Добро пожаловать в жиросброс! Снова."
	case BMI > 40:
		text = "<b>Ожирение третьей степени.</b> Если ты не тяжелоатлет, то у тебя конкретний лишний вес. Пора худеть."
	}

	msg := fmt.Sprintf(bmiMessage, m.From.FirstName, BMI, text)
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

func (b *Bot) changeRepoCommand(m *tgbotapi.Message) {
	msg := b.changeRepo(m)
	message := tgbotapi.NewMessage(m.Chat.ID, msg)
	b.BotAPI.Send(message)
}

func (b *Bot) changeRepo(m *tgbotapi.Message) string {
	if m.ReplyToMessage == nil {
		return "Используйте команду как ответ на чьё-то сообщение."
	}

	if !admins[m.From.ID] {
		return "Вы не можете использовать эту команду"
	}

	args := m.CommandArguments()
	args = strings.TrimSpace(args)

	if len(args) != 0 {
		argsArr := strings.Split(args, " ")
		amountStr := argsArr[0]
		char := argsArr[1]

		amount, err := strconv.Atoi(amountStr)
		if err != nil || amount <= 0 {
			return fmt.Sprintf("%s не является корректным числом.", amountStr)
		}

		if char != "-" && char != "+" {
			return fmt.Sprintf("%s не является корректным символом.", char)
		}

		message := tgbotapi.NewMessage(m.Chat.ID, char)
		message.ReplyToMessageID = m.ReplyToMessage.MessageID
		for i := 0; i < amount; i++ {
			res, err := b.BotAPI.Send(message)
			if err != nil {
				return fmt.Sprintf("sending msg error: %v", err)
			}

			b.BotAPI.DeleteMessage(tgbotapi.DeleteMessageConfig{
				ChatID:    res.Chat.ID,
				MessageID: res.MessageID,
			})

			time.Sleep(2 * time.Second)
		}
	}

	return "Репутация изменена."
}
