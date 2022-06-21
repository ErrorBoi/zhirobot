package whoami

import (
	"fmt"
	"math/rand"
	"time"
)

type WhoAmI struct {
}

func (w *WhoAmI) GetMyTitle() string {
	rand.Seed(time.Now().Unix())

	adjective := adjectives[rand.Intn(len(adjectives))]

	noun := nouns[rand.Intn(len(nouns))]

	return fmt.Sprintf("ты %s %s!", adjective, noun)
}

func (w *WhoAmI) GetUserTitle(nickname string) string {
	rand.Seed(time.Now().Unix())

	adjective := adjectives[rand.Intn(len(adjectives))]

	noun := nouns[rand.Intn(len(nouns))]

	return fmt.Sprintf("%s - %s %s!", nickname, adjective, noun)
}

func New() *WhoAmI {
	return &WhoAmI{}
}
