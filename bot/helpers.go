package bot

var (
	bannedUsers = []int{289675939}
	admins      = map[int]bool{128883002: true}
)

func (b *Bot) isAllowed(tgID int) bool {
	for _, item := range bannedUsers {
		if tgID == item {
			return false
		}
	}

	return true
}
