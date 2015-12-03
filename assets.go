package main

import (
	"github.com/codegangsta/cli"
)

var themeID int

func AssetCommands() cli.Command {
	return cli.Command{
		Name: "assets",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:        "theme-id",
				Usage:       "ID of the theme to use",
				Destination: &themeID,
			},
		},
		Subcommands: []cli.Command{
			cli.Command{
				Name:  "list",
				Usage: "Lists the assets in a given theme.",
			},
		},
	}
}
