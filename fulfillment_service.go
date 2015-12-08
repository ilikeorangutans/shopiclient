package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/ilikeorangutans/shopify"
	"log"
)

func FulfillmentServicesCommand() cli.Command {
	return cli.Command{
		Name:   "fulfillment-services",
		Usage:  "",
		Action: listFullfillments,
		Subcommands: []cli.Command{
			cli.Command{
				Name:   "list",
				Usage:  "Lists fulfillment services",
				Action: listFullfillments,
			},
			cli.Command{
				Name:   "create",
				Usage:  "Creates a new fulfillment service",
				Action: createFulfillmentService,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name: "name",
					},
				},
			},
		},
	}
}

func listFullfillments(c *cli.Context) {

	ffs := shopifyClient.FullfillmentServices()
	services, err := ffs.List()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Fulfillment services:", len(services))
	fmt.Printf("%-4s  %-12s  %s\n", "ID", "Name", "Handle")
	for _, s := range services {
		fmt.Printf("%-4d  %-12s  %s  %s  %s  %s\n", s.ID, s.Name, s.Handle, s.CallbackURL, s.Format, s.ProviderID)
	}

}

func createFulfillmentService(c *cli.Context) {
	ffs := shopifyClient.FullfillmentServices()
	services, err := ffs.List()
	if err != nil {
		log.Fatal(err)
	}

	_, err = ffs.Create(&shopify.FulfillmentService{
		Name: c.String("name"),
		//		CallbackURL: "http://localhost:12345/",
		Handle:      "KABAM",
		Email:       "test@test.com",
		Format:      "json",
		Credential1: "credential",
	})
	if err != nil {
		log.Fatal(err)
	}

	services, err = ffs.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range services {
		fmt.Println(s)
	}

}
