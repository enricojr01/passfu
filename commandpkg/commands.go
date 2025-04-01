package commandpkg

import (
	"github.com/urfave/cli"
)

// I'm not sure I need to be doing this in golang?
// I moved the functions out because I don't know if godocs will autogenerate
// documentation for inline functions properly (I don't think it will)
var EncryptDatabase cli.Command = cli.Command{
	Name:   "encrypt",
	Usage:  "Encrypts the contents of <infile> with a <masterpassword> and writes the encrypted content to <outfile>.",
	Action: encryptdb,
}

var DecryptDatabase cli.Command = cli.Command{
	Name:   "decrypt",
	Usage:  "Decrypts <infile> with a <masterpassword> and writes the decrypted content to <outfile>.",
	Action: decryptdb,
}

var SanityCheck cli.Command = cli.Command{
	Name:   "sanitycheck",
	Usage:  "(dev only) checks to see if encryption / decryption works in simple cases.",
	Action: sanitycheck,
}
