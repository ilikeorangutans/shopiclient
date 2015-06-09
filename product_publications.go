package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"log"
	//	"github.com/ilikeorangutans/shopify"
	//	"strings"
	"strconv"
)

func ProductPublicationsCommands() cli.Command {

	return cli.Command{
		Name:    "product-publications",
		Aliases: []string{"pp"},
		Subcommands: []cli.Command{
			cli.Command{
				Name:   "list",
				Action: ListPublications,
			},
			cli.Command{
				Name:   "create",
				Action: CreateProductPublication,
			},
			cli.Command{
				Name:   "delete",
				Action: DeleteProductPublication,
			},
		},
	}
}

func ListPublications(c *cli.Context) {
	fmt.Println("Listing publications")

	id, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		log.Fatal(err)
	}

	publications := shopifyClient.ProductPublications().List(id)

	for _, publication := range publications {
		fmt.Printf("%-5d  %5d\n", publication.ChannelId, publication.Id)

	}
}

func CreateProductPublication(c *cli.Context) {

}

func DeleteProductPublication(c *cli.Context) {}
