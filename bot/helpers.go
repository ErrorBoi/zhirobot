package bot

var (
	allowedUsers = []int{128883003}
)

func (b *Bot) isAllowed(tgID int) bool {
	for _, item := range allowedUsers {
		if tgID == item {
			return true
		}
	}

	return false
}
