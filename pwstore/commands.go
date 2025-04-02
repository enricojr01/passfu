package pwstore

import "github.com/urfave/cli"

var NewDatabase cli.Command = cli.Command{
	Name:   "newdb",
	Usage:  "Creates a new, empty, unencrypted password database.",
	Action: newdb,
}

var NewPassword cli.Command = cli.Command{
	Name:   "newpw",
	Usage:  "Adds a new password to the database",
	Action: newpw,
}

var GetPassword cli.Command = cli.Command{
	Name:   "getpw",
	Usage:  "Retrieves the record with the provided ID from the database.",
	Action: getpw,
}
