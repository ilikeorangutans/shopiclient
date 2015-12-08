package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	//	"github.com/ilikeorangutans/shopify"
	"log"
	"strings"
)

func MetafieldsCommands() cli.Command {

	return cli.Command{
		Name: "metafields",
		Subcommands: []cli.Command{
			cli.Command{
				Name:   "list",
				Action: ListMetafields,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "key, k",
						Usage: "Metafield key",
					},
					cli.StringFlag{
						Name:  "value, v",
						Usage: "Metafield value",
					},
					cli.StringFlag{
						Name:  "type, t",
						Usage: "Metafield type",
						Value: "string",
					},
				},
			},
			cli.Command{
				Name:   "create",
				Action: CreateMetafield,
			},
		},
	}
}

func CreateMetafield(c *cli.Context) {
	shopifyClient.Metafields().Create("products/34", "apes", "mykey", "myvalue", "string")
}

func ListMetafields(c *cli.Context) {
	mf := shopifyClient.Metafields()
	fields, err := mf.List("products/34")
	if err != nil {
		log.Fatal(err)
	}

	format_banner := "%-11s  %-20s  %-20s  %-10s  %-20s  %-16s\n"
	format := "%11d  %-20s  %-20s  %-10s  %-20s  %5d %-10s\n"

	fmt.Printf(format_banner, "ID", "Namespace", "Key", "Type", "Value", "On")
	for _, metafield := range fields {
		value := strings.Trim(metafield.Value, " ")
		if len(value) > 20 {
			value = strings.Trim(value[0:17], " ") + "..."
		}
		fmt.Printf(format, metafield.ID, metafield.Namespace, metafield.Key, metafield.Type, value, metafield.ResourceID, metafield.Resource)
	}
}
