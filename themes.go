package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"log"
	"os"
	"text/tabwriter"
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

	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)

	writer.Write([]byte(fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\n", "ID", "Name", "Role", "Theme Store ID", "Processing", "Previewable")))
	for i := range themesList {
		theme := themesList[i]
		writer.Write([]byte(fmt.Sprintf("%d\t%s\t%s\t%d\t%t\t%t\n", theme.ID, theme.Name, theme.Role, theme.ThemeStoreID, theme.Processing, theme.Previewable)))
	}

	writer.Flush()
}
