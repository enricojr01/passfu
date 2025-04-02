package main

import (
	"log"
	"os"
	"passfu/commandpkg"
	"passfu/pwstore"

	"github.com/rivo/tview"
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
			// yes this is messy
			// I will fix it later
			var tv *tview.TextView = tview.NewTextView()
			tv.SetBorder(true)
			tv.SetText("If you can see this...")

			var tv2 *tview.TextView = tview.NewTextView()
			tv2.SetBorder(true)
			tv2.SetText("...it's probably working")

			var tv3 *tview.TextView = tview.NewTextView()
			tv3.SetBorder(true)
			tv3.SetText("Extra Panel!")

			var grid *tview.Grid = tview.NewGrid()
			grid.SetSize(1, 3, -1, -1)
			grid.AddItem(tv, 0, 0, 1, 1, 0, 0, false)
			grid.AddItem(tv2, 0, 1, 1, 1, 0, 0, false)
			grid.AddItem(tv3, 0, 2, 1, 1, 0, 0, false)

			var napp *tview.Application = tview.NewApplication()
			napp.SetRoot(grid, true)
			// napp.SetFocus(grid)

			var err error = napp.Run()
			if err != nil {
				panic(err)
			}

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
