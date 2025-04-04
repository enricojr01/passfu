package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"passfu/commandpkg"
	"passfu/pwstore"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli"
)

type appstate struct {
	filepicker   filepicker.Model
	selectedfile string
	quitting     bool
	err          error
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (a appstate) Init() tea.Cmd {
	return a.filepicker.Init()
}

func (a appstate) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			a.quitting = true
			return a, tea.Quit
		default:
			a.quitting = false
		}
	case clearErrorMsg:
		a.err = nil
	default:
		a.quitting = false
	}

	var cmd tea.Cmd
	a.filepicker, cmd = a.filepicker.Update(msg)

	var didSelect bool
	var path string
	didSelect, path = a.filepicker.DidSelectFile(msg)
	if didSelect {
		a.selectedfile = path
	}

	didSelect, path = a.filepicker.DidSelectDisabledFile(msg)
	if didSelect {
		a.err = errors.New(path + " is not valid.")
		a.selectedfile = ""
		return a, tea.Batch(cmd, clearErrorAfter(2*time.Second))

	}

	return a, nil
}

func (a appstate) View() string {
	if a.quitting {
		return ""
	}

	var s strings.Builder
	s.WriteString("\n ")

	if a.err != nil {
		s.WriteString(a.filepicker.Styles.DisabledFile.Render(a.err.Error()))
	}

	s.WriteString("\n\n" + a.filepicker.View() + "\n")

	return s.String()
}

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
			var fp filepicker.Model = filepicker.New()
			fp.AllowedTypes = []string{}
			fp.CurrentDirectory, _ = os.Getwd()
			fmt.Println(fp.CurrentDirectory)

			var a appstate = appstate{filepicker: fp}

			var prog *tea.Program = tea.NewProgram(&a)
			var err error

			var mod tea.Model
			mod, err = prog.Run()
			if err != nil {
				return err
			}

			// gotta read up on type assertions in golang
			var mm appstate = mod.(appstate)

			fmt.Println("\n You selected: " + a.filepicker.Styles.Selected.Render(mm.selectedfile) + "\n")
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
