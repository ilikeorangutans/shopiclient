package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/ilikeorangutans/shopify"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
)

const WEBHOOK_LIST_FORMAT = "%4v  %-20s  %-6s  %-s"

func WebhookCommands() cli.Command {
	return cli.Command{
		Name:    "webhooks",
		Action:  WebhooksDefault,
		Aliases: []string{"wh"},
		Subcommands: []cli.Command{
			cli.Command{
				Name:   "list",
				Usage:  "Lists registered webhooks",
				Action: ListWebhooks,
			},
			cli.Command{
				Name:   "delete",
				Usage:  "Deletes a webhook",
				Action: DeleteWebhook,
			},
			cli.Command{
				Name:   "create",
				Usage:  "Registers a new webhook",
				Action: CreateWebhook,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name: "topic",
					},
					cli.StringFlag{
						Name: "address",
					},
					cli.StringFlag{
						Name: "format",
					},
				},
			},
			cli.Command{
				Name:   "auto-test",
				Usage:  "Automatically set up a webhook for the given topic(s) and start a server to listen",
				Action: AutoTestWebhook,
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "pretty-json",
						Usage: "Pretty print webhook JSON payload",
					},
				},
			},
		},
	}
}

func WebhooksDefault(context *cli.Context) {
	if len(context.Args()) > 0 {
		id, err := strconv.ParseInt(context.Args()[0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		webhook, _ := shopifyClient.Webhooks().Get(id)

		prettyListWebhooks(webhook)
	}
}

func prettyListWebhooks(hooks ...*shopify.Webhook) {
	if len(hooks) == 0 {
		return
	}
	fmt.Printf(WEBHOOK_LIST_FORMAT, "ID", "Topic", "Format", "Address")
	fmt.Println()
	for _, webhook := range hooks {
		fmt.Printf(WEBHOOK_LIST_FORMAT, webhook.ID, webhook.Topic, webhook.Format, webhook.Address)
		fmt.Println()
	}
}

func AutoTestWebhook(context *cli.Context) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	format := "json"

	u, err := url.Parse(fmt.Sprintf("http://%s:8080", hostname))
	if err != nil {
		log.Fatal(err)
	}

	webhooks := make([]*shopify.Webhook, 0)
	for _, topic := range context.Args() {
		webhook, err := shopifyClient.Webhooks().Create(topic, u, format)
		if err != nil {
			log.Fatal(err)
		}
		webhooks = append(webhooks, webhook)
	}

	fmt.Println("Created new webhook for automatic testing:")
	prettyListWebhooks(webhooks...)

	log.Println("Now listening for webhooks, press ^C to exit...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go interruptHandler(c, webhooks)

	responseHandler := &webhookResponseHandler{
		prettyPrint: context.Bool("pretty-json"),
	}

	http.Handle("/", responseHandler)
	http.ListenAndServe(":8080", nil)
}

func interruptHandler(c chan os.Signal, webhooks []*shopify.Webhook) {
	for _ = range c {
		for _, webhook := range webhooks {
			shopifyClient.Webhooks().Delete(webhook.ID)
		}
		os.Exit(0)
	}
}

type webhookResponseHandler struct {
	prettyPrint bool
}

func (wrh *webhookResponseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := httputil.DumpRequest(r, !wrh.prettyPrint)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
	fmt.Println()
	log.Println("--< Incoming Request >---------------------------------------------------------")
	fmt.Printf("%s", b)

	if wrh.prettyPrint {
		decoder := json.NewDecoder(r.Body)

		var v interface{}
		decoder.Decode(&v)

		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s", b)
	}

	fmt.Println()
}

func CreateWebhook(context *cli.Context) {
	u, err := url.Parse(context.String("address"))
	if err != nil {
		log.Fatal(err)
	}

	format := strings.Trim(context.String("format"), " ")
	if format != "json" && format != "xml" {
		log.Fatalf("Invalid format %s, expected either \"json\" or \"xml\"!", format)
	}

	topic := context.String("topic")

	webhook, err := shopifyClient.Webhooks().Create(topic, u, format)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created new webhook:")
	prettyListWebhooks(webhook)
}

func ListWebhooks(context *cli.Context) {
	webhooks := shopifyClient.Webhooks()
	hooks, err := webhooks.List()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Registered webhooks: %d (you only see webhooks registered with the current credentials)\n", len(hooks))
	prettyListWebhooks(hooks...)
}

func DeleteWebhook(context *cli.Context) {
	webhooks := shopifyClient.Webhooks()

	id, err := strconv.ParseInt(context.Args()[0], 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	webhooks.Delete(id)
}
