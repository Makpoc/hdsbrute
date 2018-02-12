package coffee

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/makpoc/hdsbrute"
)

// CoffeeCommand is dummy command to share a coffe recipe
var CoffeeCommand = hdsbrute.Command{
	Cmd:     "coffee",
	HelpStr: "Coffee is Love, Coffe is Life",
	Init: func() error {
		fmt.Println("Coffee brewing!")
		return nil
	},
	Exec: coffeeHandler,
}

// Handle responds to .coffee command with some recepies
func coffeeHandler(s *discordgo.Session, m *discordgo.MessageCreate, query []string) {
	_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(`Here's the recipe you asked for %s:
    **1.** Brew a large espresso :coffee:.
    **2.** Fill a cocktail shaker half full with ice cubes :cocktail:.
    **3.** Add to the ice the brewed espresso.
    **4.** 3 tablespoons vodka.
    **5.** 3 tablespoons Kahlúa (coffee liqueur)
    **6.** 1/4 teaspoon sugar.
    **7.** Shake until foamy, about 30 seconds; strain into a martini glass.
    **8.** Give to someone else and get a :beer: or 3
    **9.** Enjoy :beers:`, m.Author.Mention()))
	if err != nil {
		fmt.Println(err)
	}
}
