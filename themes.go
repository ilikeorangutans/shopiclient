package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"log"
)

func ThemeCommands() cli.Command {
	return cli.Command{
		Name: "themes",
		Subcommands: []cli.Command{

			cli.Command{
				Name:   "list",
				Usage:  "Lists available themes",
				Action: ListThemes,
			},
		},
	}
}

func ListThemes(context *cli.Context) {
	themes := shopifyClient.Themes()

	themesList, err := themes.List()
	if err != nil {
		log.Fatal(err)
	}

	for i := range themesList {
		theme := themesList[i]
		fmt.Printf("%v\n", theme)

	}

}
