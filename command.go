package hdsbrute

import "github.com/bwmarrin/discordgo"

// HandlerFn is the function to handle commands
type HandlerFn func(*discordgo.Session, *discordgo.MessageCreate, []string)

// Command represents a command
type Command struct {
	Cmd     string
	Args    []string
	HelpStr string
	HelpFn  func(*Brute, *discordgo.Session, *discordgo.MessageCreate)
	Init    func() error
	Exec    HandlerFn
}
