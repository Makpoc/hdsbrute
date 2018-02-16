package main

import (
	"log"
	"os"

	"github.com/makpoc/hdsbrute"
	"github.com/makpoc/hdsbrute/commands/coffee"
	"github.com/makpoc/hdsbrute/commands/frosty"
	"github.com/makpoc/hdsbrute/commands/gsheet"
	"github.com/makpoc/hdsbrute/commands/sheet"
	"github.com/makpoc/hdsbrute/commands/wsmap"
)

func main() {
	brute, err := hdsbrute.New(hdsbrute.GetEnvPropOrDefault("BOT_TOKEN", ""))

	if err != nil {
		printAndExit(err)
	}

	brute.AddCommand(frosty.FrostyCommand)
	brute.AddCommand(coffee.CoffeeCommand)
	brute.AddCommand(sheet.SheetCommand)
	brute.AddCommand(wsmap.WsCommand)
	brute.AddCommand(gsheet.TimeZoneCommand)

	log.Println("Bot is running")

	err = brute.Start()
	if err != nil {
		printAndExit(err)
	}
	<-make(chan struct{})
}

func printAndExit(err error) {
	log.Printf("%v\n", err)
	os.Exit(1)
}
