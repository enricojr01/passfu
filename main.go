package main

import (
	"fmt"
	"passfu/easycipher"

	"gorm.io/gorm"
)

type Record struct {
	gorm.Model
	Name     string
	Username string
	Password string
}

func main() {
	var password string = "Super secret password"
	var plaintext string = "Have you heard of the critically acclaimed MMORPG Final Fantasy XIV?"
	var ciphertext []byte
	var err error

	var easyCipher easycipher.EasyCipher
	easyCipher, err = easycipher.NewEC(password)
	if err != nil {
		panic(err)
	}

	fmt.Println("plaintext: ", plaintext)
	fmt.Println("key: ", easyCipher.Key)
	fmt.Println("salt: ", easyCipher.Salt)
	fmt.Println("iv: ", easyCipher.Iv)
	easyCipher.Encrypt(plaintext)
	fmt.Println(easyCipher.Ciphertext)
	ciphertext = easyCipher.ExportCiphertext()
	fmt.Println(ciphertext)
	fmt.Println()
	fmt.Println()

	// now the other way!

	var easy2 easycipher.EasyCipher = easycipher.NewECFromCiphertext(ciphertext, password)
	fmt.Println("key: ", easy2.Key)
	fmt.Println("salt: ", easy2.Salt)
	fmt.Println("iv: ", easy2.Iv)
	var plaintext2 []byte = easy2.Decrypt()

	fmt.Println("plaintext: ", string(plaintext2))
}
