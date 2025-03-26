package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var NewDatabase cli.Command = cli.Command{
	Name:  "newdb",
	Usage: "Creates a new, empty, unencrypted password database.",
	Action: func(ctx *cli.Context) error {
		var currenttime time.Time = time.Now()
		var now int64 = currenttime.Unix()
		var filename string = fmt.Sprintf("unencrypted-%d.db", now)

		var sqlitedb gorm.Dialector = sqlite.Open(filename)

		var db *gorm.DB
		var err error
		db, err = gorm.Open(sqlitedb)
		if err != nil {
			return err
		}

		db.AutoMigrate(&Record{})

		return nil
	},
}

var EncryptDatabase cli.Command = cli.Command{
	Name:  "encryptdb",
	Usage: "Encrypts a database file with a master password.",
	Action: func(ctx *cli.Context) error {
		var args cli.Args = ctx.Args()

		var file string = args[0]
		var outfile string = args[1]
		var masterpass string = args[2]

		var ec EasyCipher
		var data []byte
		var encryptedata []byte
		var err error

		if len(args) != 3 {
			return errors.New("usage: encryptdb <infile> <outfile> <masterpassword>")
		}

		data, err = os.ReadFile(file)
		if err != nil {
			return err
		}

		ec, err = NewEC(masterpass)
		if err != nil {
			return err
		}
		ec.Encrypt(data)

		encryptedata = ec.ExportCiphertext()

		err = os.WriteFile(outfile, encryptedata, 0644)
		if err != nil {
			return err
		}

		fmt.Println("If you see this, the follwing file has been written: ", outfile)

		return nil
	},
}
