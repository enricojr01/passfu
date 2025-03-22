package main

import (
	"fmt"

	"github.com/urfave/cli"
)

var NewDatabase cli.Command = cli.Command{
	Name:  "newdb",
	Usage: "Creates a new, empty, unencrypted password database.",
	Action: func(ctx *cli.Context) error {
		fmt.Println("Pretend this command makes a new database.")
		return nil
	},
}
