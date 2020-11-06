package bot

var (
	allowedUsers = map[int]bool{128883002: true, 207660019: true}
)

func (b *Bot) isAllowed(tgID int) bool {
	return allowedUsers[tgID]
}
