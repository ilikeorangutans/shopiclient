package main

import (
	"github.com/codegangsta/cli"
	"github.com/ilikeorangutans/shopify"
	"log"
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
			Name:  "dump-urls",
			Usage: "Show URLs being requested",
		},
		cli.BoolFlag{
			Name:  "dump-requests",
			Usage: "Dump full HTTP requests as they are issued",
		},
		cli.BoolFlag{
			Name:  "dump-responses",
			Usage: "Dump full HTTP responses as they are received",
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
		ProductCommands(),
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
	if context.Command.Name == "help" {
		return nil
	}
	//(context.String("host"), context.String("user"), context.String("password"))
	settings := shopify.ClientSettings{Host: context.String("host"), Username: context.String("user"), Password: context.String("password")}
	shopifyClient = shopify.NewClientWithSettings(settings)
	err := shopifyClient.Connect()
	if err != nil {
		log.Fatal(err)
	}
	shopifyClient.Verbose = context.IsSet("verbose")
	return nil
}
