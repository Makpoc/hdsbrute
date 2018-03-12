package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/makpoc/hdsbrute"
	"github.com/makpoc/hdsbrute/commands/coffee"
	"github.com/makpoc/hdsbrute/commands/frosty"
	"github.com/makpoc/hdsbrute/commands/sheet"
	"github.com/makpoc/hdsbrute/commands/tz"
	"github.com/makpoc/hdsbrute/commands/userinfo"
	"github.com/makpoc/hdsbrute/commands/wsmap"
)

func main() {
	brute, err := hdsbrute.New(hdsbrute.GetEnvPropOrDefault("BOT_TOKEN", ""))

	if err != nil {
		printAndExit(err)
	}

	brute.AddCommand(frosty.FrostyCommand)
	brute.AddCommand(coffee.CoffeeCommand)
	brute.AddCommand(sheet.SGSheetCommand)
	brute.AddCommand(sheet.CTSheetCommand)
	brute.AddCommand(wsmap.WsCommand)
	brute.AddCommand(tz.TimeZoneCommand)
	brute.AddCommand(userinfo.UserInfoCommand)

	err = brute.Start()
	if err != nil {
		printAndExit(err)
	}
	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	err = brute.Close()
	if err != nil {
		printAndExit(err)
	}
}

func printAndExit(err error) {
	log.Printf("%v\n", err)
	os.Exit(1)
}
