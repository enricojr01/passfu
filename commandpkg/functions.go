package commandpkg

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"passfu/easycipher"

	"github.com/urfave/cli"
)

// This function powers the "sanitycheck" command, which is a dev command used
// to make sure that the whole encrypt/decrypt workflow is working correctly.
// It was used during development to refine the design of the easycipher package
// and remains in-place as a tool for future debugging.
func sanitycheck(ctx *cli.Context) error {
	// var args cli.Args = ctx.Args()

	var instring string = "Do you know about the award-winning MMO Final Fantasy XIV?"
	var password string = "testpass1"
	var err error

	var ec easycipher.EasyCipher
	ec, err = easycipher.New(password, []byte(instring))
	if err != nil {
		return err
	}

	ec.Encrypt()
	fmt.Println("ec.Salt: ", ec.Salt)
	fmt.Println("ec.Iv: ", ec.Iv)
	fmt.Println("ec.Key: ", ec.Key)
	fmt.Println("ec.Ciphertext: ", ec.Ciphertext)

	var ciphertext []byte = ec.Ciphertext
	var ec2 easycipher.EasyCipher
	ec2, err = easycipher.NewFromCiphertext(password, ciphertext)
	if err != nil {
		return err
	}

	fmt.Println("ec2.Salt: ", ec2.Salt)
	fmt.Println("ec2.Iv: ", ec2.Iv)
	fmt.Println("ec2.Key: ", ec2.Key)
	fmt.Println("ec2.Ciphertext: ", ec2.Ciphertext)

	fmt.Println("same salt? ", bytes.Equal(ec.Salt, ec2.Salt))
	fmt.Println("same IV? ", bytes.Equal(ec.Iv, ec2.Iv))
	fmt.Println("same key? ", bytes.Equal(ec.Key, ec2.Key))
	fmt.Println("same plaintext? ", bytes.Equal(ec.Plaintext, ec2.Plaintext))

	ec2.Decrypt()
	fmt.Println("ec2.Plaintext: ", string(ec2.Plaintext))
	return nil
}

func encryptdb(ctx *cli.Context) error {
	var args cli.Args = ctx.Args()

	if len(args) != 3 {
		return errors.New("usage: encryptdb <infile> <outfile> <masterpassword>")
	}

	var file string = args[0]
	var outfile string = args[1]
	var masterpass string = args[2]

	var ec easycipher.EasyCipher
	var data []byte
	var err error

	data, err = os.ReadFile(file)
	if err != nil {
		return err
	}

	ec, err = easycipher.New(masterpass, data)
	if err != nil {
		return err
	}
	fmt.Println("ec.Salt: ", ec.Salt)
	fmt.Println("ec.Iv: ", ec.Iv)
	ec.Encrypt()

	err = os.WriteFile(outfile, ec.Ciphertext, 0644)
	if err != nil {
		return err
	}

	fmt.Println("If you see this, the follwing file has been written: ", outfile)

	return nil
}

func decryptdb(ctx *cli.Context) error {
	var args cli.Args = ctx.Args()

	if len(args) != 3 {
		return errors.New("usage: decrypt <infile> <outfile> <masterpassword>")
	}

	var infile string = args[0]
	var outfile string = args[1]
	var masterpass string = args[2]

	var data []byte
	var err error
	data, err = os.ReadFile(infile)
	if err != nil {
		return err
	}

	var ec easycipher.EasyCipher
	ec, err = easycipher.NewFromCiphertext(masterpass, data)
	if err != nil {
		return err
	}
	fmt.Println("ec.Salt: ", ec.Salt)
	fmt.Println("ec.Iv: ", ec.Iv)
	ec.Decrypt()

	err = os.WriteFile(outfile, ec.Plaintext, 0644)
	if err != nil {
		return err
	}

	fmt.Println("If you see this, the follwing file has been written: ", outfile)

	return nil
}
