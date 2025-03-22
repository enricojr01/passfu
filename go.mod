module passfu

go 1.24.1

require gorm.io/gorm v1.25.12

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/text v0.23.0 // indirect
	passfu/easycipher v0.0.0-00010101000000-000000000000 // indirect
)

replace passfu/easycipher => ./easycipher
