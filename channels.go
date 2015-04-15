package main

import (
	"fmt"
	"github.com/codegangsta/cli"
)

func ChannelCommands() cli.Command {
	return cli.Command{
		Name:   "channels",
		Action: ChannelsDefault,
	}
}

func ChannelsDefault(context *cli.Context) {
	channels := shopifyClient.Channels()
	ch := channels.List()

	fmt.Printf("Found %d channels.\n", len(ch))
	fmt.Printf("%-5s %-10s\n", "ID", "Provider ID")
	for _, c := range ch {
		fmt.Printf("%-5d %-10d\n", c.Id, c.ProviderId)
	}
}
