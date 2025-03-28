package main

import (
	"fmt"
	"log"
	"os"

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
	commands = append(commands, NewDatabase)
	commands = append(commands, EncryptDatabase)
	commands = append(commands, DecryptDatabase)
	commands = append(commands, SanityCheck)

	var app *cli.App = &cli.App{
		Name:  "passfu",
		Usage: "A password manager for the command line.",
		Action: func(*cli.Context) error {
			fmt.Println("Test! If you see this it's working.")
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
