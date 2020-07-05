package bot

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) GetWeightKeyboard(tgID, page int, firstName string, last bool) *tgbotapi.InlineKeyboardMarkup {
	var buttons []tgbotapi.InlineKeyboardButton

	// Add "back" button for all pages except first one
	if page != 0 {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("◀️Назад", fmt.Sprintf("getWeight-%d-%d-%s", tgID, page-1, firstName)))
	}

	buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(page+1), "none"))

	// Add "forward" button for all pages except last one
	if !last {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("▶ Вперёд", fmt.Sprintf("getWeight-%d-%d-%s", tgID, page+1, firstName)))
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(buttons...),
		)

	return &keyboard
}
