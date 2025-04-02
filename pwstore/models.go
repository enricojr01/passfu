package pwstore

import (
	"gorm.io/gorm"
)

type Record struct {
	gorm.Model
	Name     string
	Username string
	Password string
	Notes    string
}

func New(name string, username string, password string, notes string) Record {
	var rec Record = Record{
		Name:     name,
		Username: username,
		Password: password,
		Notes:    notes,
	}
	return rec
}
