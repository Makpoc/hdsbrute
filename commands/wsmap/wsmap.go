package wsmap

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/makpoc/hdsbrute"
)

type mapCommand struct {
	args    []string
	author  string
	message []string
}

var backendURL string
var backendSecret string

const backtick = "&bt;"
const cmd = "map"

// WsCommand ...
var WsCommand = hdsbrute.Command{
	Cmd:      cmd,
	HelpFunc: helpFunc,
	Init: func(b *hdsbrute.Brute) error {
		backendSecret = b.Config.Secret
		backendURL = b.Config.BackendURL

		log.Println("Map initialized")
		return nil
	},
	Exec: handleFunc,
}

// helpFunc is the function called to display help/usage info
func helpFunc(b *hdsbrute.Brute, s *discordgo.Session, m *discordgo.MessageCreate) {
	helpMessage := buildHelpString(fmt.Sprintf("%s%s", b.Prefix, cmd))
	s.ChannelMessageSend(m.ChannelID, helpMessage)
}

// buildHelpString builds the help and usage string for the map
func buildHelpString(cmdWithPrefix string) string {
	header := fmt.Sprintf("`%s` displays the map for the current WS with coordinates overlayed on top and optional sector highlighting.\n\n**Usage**: `%s [[color] [coord...]...] message`\n\n**Examples**:", cmdWithPrefix, cmdWithPrefix)
	var subCommands []string

	subCommands = append(subCommands, fmt.Sprintf("`%s` - Displays just the map with the coordinates", cmdWithPrefix))
	subCommands = append(subCommands, fmt.Sprintf("`%s [coord...]` - Highlights the sector(s) specified by the provided coordinates. e.g. `%s a1 b2 c3`", cmdWithPrefix, cmdWithPrefix))
	subCommands = append(subCommands, fmt.Sprintf("`%s [[color] [coord...]...]` - Highlights the provided coordinates with the color that comes before them. E.g. `%s orange a3 green b4 b5 pink d3`. Currently supported colors are: **orange**, **yellow**, **green**, **pink**.", cmdWithPrefix, cmdWithPrefix))
	subCommands = append(subCommands, fmt.Sprintf("`%s [color|coords...] [message]` - Same as above but also adds message with details. E.g. `%s orange d3 pink b2 b3 BS defend at d3. Miners void b2 b3`.", cmdWithPrefix, cmdWithPrefix))

	finalMessage := fmt.Sprintf("%s\n%s", header, strings.Join(subCommands, "\n"))

	return strings.TrimRight(finalMessage, "\n")
}

// handlerFunc answers calls to map and map [coord|color] message
func handleFunc(b *hdsbrute.Brute, s *discordgo.Session, m *discordgo.MessageCreate, query []string) {
	mCommand := parseMapCommand(query)
	mCommand.author = m.Author.Username

	url := fmt.Sprintf("%s/api/v1/map?secret=%s", backendURL, backendSecret)
	if len(mCommand.args) > 0 {
		url = fmt.Sprintf("%s&coords=%s", url, strings.Join(mCommand.args, ","))
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to get map. Error was: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":flushed: Failed to get map"))
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code: %v\n", resp.Status)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":flushed: Failed to get map"))
		return
	}

	err = sendDiscordResponse(s, m, resp, mCommand)
	if err != nil {
		log.Println("Something went wrong while sending Discord response", err)
		return
	}

	err = s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		log.Println("Failed to delete trigger message", err)
	}
}

// sendDiscordResponse sends the response from the backend to the discord channel it got the trigger from. It will also add a message to the file in that response, containing the author of the trigger and will delete the original message.
func sendDiscordResponse(s *discordgo.Session, m *discordgo.MessageCreate, resp *http.Response, mCommand mapCommand) error {
	respContentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(respContentType, "image/") {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(":thinking: Suspecious content type: %s!", respContentType))
		return fmt.Errorf("invalid map format")
	}

	var err error
	if len(mCommand.message) > 0 {
		message := fmt.Sprintf("**%s**: %s", mCommand.author, strings.Join(mCommand.message, " "))
		_, err = s.ChannelFileSendWithMessage(m.ChannelID, message, "map.jpeg", resp.Body)
	} else {
		_, err = s.ChannelFileSendWithMessage(m.ChannelID, fmt.Sprintf("**%s** asked for: ", mCommand.author), "map.jpeg", resp.Body)
	}

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, ":flushed: Failed to fullfil your desire")
		return err
	}

	return nil
}

// parseMapCommand parses the given command into a map command struct
func parseMapCommand(query []string) mapCommand {
	var mCommand mapCommand
	if len(query) > 0 {
		var mIndex int
		var w string
		for _, w = range query {
			if !isValidArgument(w) {
				break
			}
			mCommand.args = append(mCommand.args, w)
			mIndex++
		}
		mCommand.message = query[mIndex:]
	}
	return mCommand
}

// isValidArgument checks if the provided string is a valid argument (coordinate or color)
func isValidArgument(arg string) bool {
	directions := []string{
		"a1", "a2", "a3", "a4",
		"b1", "b2", "b3", "b4", "b5",
		"c1", "c2", "c3", "c4", "c5", "c6",
		"d1", "d2", "d3", "d4", "d5", "d6", "d7",
		"e1", "e2", "e3", "e4", "e5", "e6",
		"f1", "f2", "f3", "f4", "f5",
		"g1", "g2", "g3", "g4",
	}

	colors := []string{
		"green", "orange", "pink", "yellow", "red", "warn",
	}

	arg = strings.ToLower(arg)

	return hdsbrute.Contains(directions, arg) || hdsbrute.Contains(colors, arg)
}
