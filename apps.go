package main

import (
	"fmt"
	"github.com/codegangsta/cli"
)

func foo(context *cli.Context) {
	aps := shopifyClient.Apps().List()

	for _, ap := range aps {
		fmt.Printf("%d \n", ap.Id)
	}
}
