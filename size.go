package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/digitalocean/godo"

	"golang.org/x/oauth2"
)

var SizeCommand = cli.Command{
	Name:   "size",
	Usage:  "Size commands.",
	Action: sizeList,
	Subcommands: []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"s"},
			Usage:   "List all sizes.",
			Action:  sizeList,
		},
	},
}

func sizeList(ctx *cli.Context) {
	if ctx.BoolT("help") == true {
		cli.ShowAppHelp(ctx)
		os.Exit(1)
	}

	tokenSource := &TokenSource{
		AccessToken: APIKey,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)

	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 50, // Not likely to have more than 50 sizes soon
	}
	sizeList, _, err := client.Sizes.List(opt)
	if err != nil {
		fmt.Printf("Unable to list Sizes: %s\n", err)
		os.Exit(1)
	}

	cliOut := NewCLIOutput()
	defer cliOut.Flush()
	cliOut.Header("Slug", "Memory", "VCPUs", "Disk", "Transfer", "Price Monthly", "Price Hourly")
	for _, size := range sizeList {
		cliOut.Writeln("%s\t%dMB\t%d\t%dGB\t%d\t$%.0f\t$%.5f\n",
			size.Slug, size.Memory, size.Vcpus, size.Disk, size.Transfer, size.PriceMonthly, size.PriceHourly)
	}
}
