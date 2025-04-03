package main

import (
	"fmt"
	"log"
	"os"
	"passfu/commandpkg"
	"passfu/pwstore"

	"github.com/urfave/cli"
)

func main() {
	var me cli.Author = cli.Author{
		Name:  "Enrico Tuvera Jr",
		Email: "test@gmail.com",
	}

	var authors []cli.Author
	authors = append(authors, me)

	var commands []cli.Command
	commands = append(commands, pwstore.NewDatabase)
	commands = append(commands, pwstore.NewPassword)
	commands = append(commands, pwstore.GetPassword)
	commands = append(commands, commandpkg.EncryptDatabase)
	commands = append(commands, commandpkg.DecryptDatabase)
	commands = append(commands, commandpkg.SanityCheck)

	var app *cli.App = &cli.App{
		Name:  "passfu",
		Usage: "A password manager for the command line.",
		Action: func(*cli.Context) error {
			fmt.Println("todo, will add stuff here later.")
			return nil
		},
		Authors:  authors,
		Commands: commands,
	}

	var err error = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
