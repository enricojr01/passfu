package pwstore

import (
	"gorm.io/gorm"
)

type Record struct {
	gorm.Model
	Name     string
	Username string
	Password string
}
