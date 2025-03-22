package main

import (
	"fmt"
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

// var EncryptDatabase cli.Command = cli.Command{
// 	Name: "encryptdb"
// }
