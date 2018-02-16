package hdsbrute

import "github.com/bwmarrin/discordgo"

// HandlerFn is the function to handle commands
type HandlerFn func(*discordgo.Session, *discordgo.MessageCreate, []string)

// Command represents a command
type Command struct {
	Cmd      string
	Args     []string
	HelpStr  string
	HelpFunc func(*Brute, *discordgo.Session, *discordgo.MessageCreate)
	Init     func(*Brute) error
	Exec     HandlerFn
}
