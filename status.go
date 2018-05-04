package hdsbrute

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

// ready will be called when the bot receives the "ready" event from Discord.
func ready(s *discordgo.Session, event *discordgo.Ready) {
	go func() {
		var statuses = []string{
			"Hades' Star",
			"RS with bots",
			"RS, killing cerb",
			"cards with colossi",
			"shipments delivery",
			"shipments delivery",
			"shipments delivery",
			"but low on hydro",
			"with TM variations",
			"with sand on Mars",
		}
		rand.Seed(time.Now().Unix())
		for {
			err := s.UpdateStatus(0, statuses[rand.Intn(len(statuses))])
			if err != nil {
				fmt.Printf("%#v\n", err)
			}
			time.Sleep(20 * time.Minute)
		}
	}()
}
