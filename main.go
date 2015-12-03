package main

import (
	"github.com/codegangsta/cli"
	"github.com/ilikeorangutans/shopify"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "shopiclient"
	app.Usage = "Shopify CLI API client"
	app.Author = "Jakob Külzer (jakob.kulzer@shopify.com)"
	app.Before = SetupClient
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Show HTTP responses",
		},
		cli.StringFlag{
			Name:   "user",
			EnvVar: "SHOPIFY_USER",
		},
		cli.StringFlag{
			Name:   "password",
			EnvVar: "SHOPIFY_PASSWORD",
		},
		cli.StringFlag{
			Name:   "host",
			EnvVar: "SHOPIFY_HOST",
		},
	}
	app.Commands = []cli.Command{
		WebhookCommands(),
		MetafieldsCommands(),
		FulfillmentServicesCommand(),
		OrdersCommands(),
		ThemeCommands(),
		AssetCommands(),
	}

	app.Run(os.Args)
}

var shopifyClient *shopify.Client

func SetupClient(context *cli.Context) error {
	shopifyClient = shopify.NewClient(context.String("host"), context.String("user"), context.String("password"))
	shopifyClient.Verbose = context.IsSet("verbose")
	return nil
}
