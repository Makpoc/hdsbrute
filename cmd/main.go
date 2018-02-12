package main

import (
	"fmt"
	"os"

	"github.com/makpoc/hdsbrute"
	"github.com/makpoc/hdsbrute/commands/coffee"
	"github.com/makpoc/hdsbrute/commands/frosty"
	"github.com/makpoc/hdsbrute/commands/sheet"
	"github.com/makpoc/hdsbrute/commands/wsmap"
)

var token string

var botID string
var botPrefix = "!"

var (
	dbUser string
	dbPass string
	dbName string
)

func main() {
	initEnv()

	brute, err := hdsbrute.New(botPrefix, token)

	if err != nil {
		printAndExit(err)
	}

	fmt.Println(brute.BotID)

	brute.AddCommand(frosty.FrostyCommand)
	brute.AddCommand(coffee.CoffeeCommand)
	brute.AddCommand(sheet.SheetCommand)
	brute.AddCommand(wsmap.WsCommand)

	fmt.Println("Bot is running")

	err = brute.Start()
	if err != nil {
		printAndExit(err)
	}
	<-make(chan struct{})
}

// initEnv initializes the application from the environment
func initEnv() {
	token = hdsbrute.GetEnvPropOrDefault("BOT_TOKEN", "")

	dbPass = hdsbrute.GetEnvPropOrDefault("dbPass", "")
	dbName = hdsbrute.GetEnvPropOrDefault("dbName", "")
	dbUser = hdsbrute.GetEnvPropOrDefault("dbUser", "")

}
func printAndExit(err error) {
	fmt.Printf("%v\n", err)
	os.Exit(1)
}
