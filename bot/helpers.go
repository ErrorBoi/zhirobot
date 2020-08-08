package bot

var (
	bannedUsers = []int{289675939}
)

func (b *Bot) isAllowed(tgID int) bool {
	for _, item := range bannedUsers {
		if tgID == item {
			return false
		}
	}

	return true
}
