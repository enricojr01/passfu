package pwstore

import (
	"errors"
	"fmt"
	"time"

	"github.com/urfave/cli"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newdb(ctx *cli.Context) error {
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
	fmt.Println("If you see this, the follwing file has been written: ", filename)
	return nil
}

func newpw(ctx *cli.Context) error {
	var args cli.Args = ctx.Args()
	if len(args) != 5 {
		return errors.New("usage: newpw <filename> <recordname> <username> <pw> <note>")
	}

	var filename string = args[0]
	var recordname string = args[1]
	var username string = args[2]
	var password string = args[3]
	var note string = args[4]
	var pwrecord Record = Record{
		Name:     recordname,
		Username: username,
		Password: password,
		Notes:    note,
	}

	fmt.Println("record: ", pwrecord)

	var sqlitedb gorm.Dialector = sqlite.Open(filename)
	var db *gorm.DB
	var err error

	db, err = gorm.Open(sqlitedb)
	if err != nil {
		return err
	}

	db.Create(&pwrecord)
	fmt.Println("A new record was created: ", pwrecord)

	return nil
}

func getpw(ctx *cli.Context) error {
	var args cli.Args = ctx.Args()
	if len(args) != 2 {
		return errors.New("usage: getpw filename recordID")
	}

	var filename string = args[0]
	var recordid string = args[1]

	var sqlitedb gorm.Dialector = sqlite.Open(filename)
	var db *gorm.DB
	var err error

	db, err = gorm.Open(sqlitedb)
	if err != nil {
		return err
	}

	var rec Record
	db.First(&rec, recordid)

	fmt.Println("record: ", rec)

	return nil
}
