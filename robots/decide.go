package robots

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type DecideBot struct {
}

func init() {
	d := &DecideBot{}
	RegisterRobot("decide", d)
}

func (d DecideBot) Run(p *Payload) (slashCommandImmediateReturn string) {
	go d.DeferredAction(p)
	text := strings.TrimSpace(p.Text)
	if text == "" {
		return "I need something to decide on!"
	} else {
		return ""
	}
}

func (d DecideBot) DeferredAction(p *Payload) {
	text := strings.TrimSpace(p.Text)
	if text != "" {
		split := strings.Split(text, ",")

		for idx, val := range split{
			split[idx] = strings.TrimSpace(val);
		}

		response := &IncomingWebhook{
			Channel:     p.ChannelID,
			Username:    "Fate Bot",
			Text:        fmt.Sprintf("@%s: Deciding between: (%s)", p.UserName, strings.Join(split, ", ")),
			IconEmoji:   ":ghost:",
			UnfurlLinks: true,
			Parse:       ParseStyleFull,
		}
		response2 := &IncomingWebhook{
			Channel:     p.ChannelID,
			Username:    "Fate Bot",
			Text:        fmt.Sprintf("@%s: Decided on: %s", p.UserName, Decide(split)),
			IconEmoji:   ":ghost:",
			UnfurlLinks: true,
			Parse:       ParseStyleFull,
		}

		response.Send()
		response2.Send()
	}
}

func (d DecideBot) Description() (description string) {
	return "Decides your fate!\n\tUsage: /decide Life, Death ...\n\tExpected Response: Deciding on (Life, Death, ...)\n\tDecided on Life!"
}

func Decide(Fates []string) (result string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := len(Fates)
	if n > 0 {
		return Fates[r.Intn(n)]
	} else {
		return fmt.Sprintf("Error")
	}
}
