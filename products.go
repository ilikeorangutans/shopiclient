package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"log"
	"os"
	"text/tabwriter"
)

func ProductCommands() cli.Command {
	return cli.Command{
		Name:  "products",
		Usage: "Products",
		Subcommands: []cli.Command{
			cli.Command{
				Name:   "list",
				Action: ListProducts,
			},
		},
	}
}

func ListProducts(context *cli.Context) {
	prods := shopifyClient.Products()
	products, err := prods.List()
	if err != nil {
		log.Fatal(err)
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	writer.Write([]byte(fmt.Sprintf("%s\t%s\n", "ID", "Title")))
	for i := range products {
		product := products[i]
		writer.Write([]byte(fmt.Sprintf("%d\t%s\n", product.ID, product.Title)))
	}
	writer.Flush()
}
