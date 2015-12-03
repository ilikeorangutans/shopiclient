package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"log"
	"os"
	"strconv"
	"text/tabwriter"
)

func OrdersCommands() cli.Command {
	return cli.Command{
		Name:    "orders",
		Aliases: []string{"o"},
		Subcommands: []cli.Command{
			cli.Command{
				Name:   "list",
				Usage:  "Lists orders",
				Action: ListOrders,
			},
			cli.Command{
				Name:    "transactions",
				Aliases: []string{"t"},
				Action:  ListOrderTransactions,
			},
		},
	}
}

func ListOrderTransactions(c *cli.Context) {

	id, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		log.Fatal(err)
	}

	order, err := shopifyClient.Orders().Get(id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Listing transactions for order #%d (%d)\n", order.Number, order.ID)

	transactions := shopifyClient.Transactions().List(order)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 5, 4, 2, ' ', 0)
	fmt.Fprintf(w, "ID\tAmount\tCurrency\tStatus\tTest\n")
	for _, transaction := range transactions {
		fmt.Fprintf(w, "%d\t%.2f\t%s\t%s\t%v\n", transaction.ID, transaction.Amount, transaction.Currency, transaction.Status, transaction.Test)
	}
	w.Flush()
}

func ListOrders(c *cli.Context) {
	orders := shopifyClient.Orders().List()

	fmt.Printf("Listing orders (%d): \n", len(orders))

	for _, o := range orders {
		fmt.Printf("%-4d %-4d %6s\n", o.ID, o.Number, o.Name)
	}
}
