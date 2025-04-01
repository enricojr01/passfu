package pwstore

import (
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
