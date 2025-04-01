package pwstore

import "github.com/urfave/cli"

var NewDatabase cli.Command = cli.Command{
	Name:   "newdb",
	Usage:  "Creates a new, empty, unencrypted password database.",
	Action: newdb,
}
