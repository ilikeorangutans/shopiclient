package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"log"
)

func ListApps(context *cli.Context) {
	apps, err := shopifyClient.Apps().List()
	if err != nil {
		log.Fatal(err)
	}

	for _, ap := range apps {
		fmt.Printf("%d \n", ap.ID)
	}
}
