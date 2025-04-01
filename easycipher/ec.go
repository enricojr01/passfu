package easycipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

// EasyCipher is a utility class designed to make encrypting / decrypting
// easier.
// The encryption method used in this module is almost identical to the one
// described in the Bitwarden whitepaper:
// A "master password" is fed to PBKDF2 to generate a key that is used
// to AES256 encrypt any arbitrary stream of bytes.
// This functionality is meant to be used by passfu to encrypt / decrypt an
// sqlite database.
// One should not instantiate EasyCipher directly, and instead use the New()
// and NewFromCiphertext() helper functions to generate them.
type EasyCipher struct {
	Password   string
	Key        []byte
	Salt       []byte
	Iv         []byte
	Ciphertext []byte
	Plaintext  []byte
}

// This function generates a new EasyCipher{} from a "master password" and
// a stream of bytes. It is meant to be used on data that has yet to be
// encrypted.
// example usage:
// var ec easycipher.EasyCipher = New("secretpw", []byte("We attack at dawn."))
// ec.Encrypt()
// // ec.Ciphertext can be written to file or used wherever
// fmt.Println(ec.Ciphertext)
func New(password string, plaintext []byte) (EasyCipher, error) {
	var newSalt []byte
	var newIv []byte
	var newKey []byte
	var err error

	newSalt, err = gimmeSalt()
	if err != nil {
		return EasyCipher{}, err
	}

	newKey, err = gimmeKey(password, newSalt)
	if err != nil {
		return EasyCipher{}, err
	}

	newIv, err = gimmeIV()
	if err != nil {
		return EasyCipher{}, err
	}

	var ec EasyCipher = EasyCipher{
		Password:   password,
		Key:        newKey,
		Salt:       newSalt,
		Iv:         newIv,
		Ciphertext: nil,
		Plaintext:  plaintext,
	}

	return ec, nil
}

// This function generates an EasyCipher{} from an already-encrypted stream
// of bytes, and a password.
func NewFromCiphertext(password string, ciphertext []byte) (EasyCipher, error) {
	var salt []byte = ciphertext[:16]
	var iv []byte = ciphertext[16 : 16+12]
	var cleanciphertext []byte = ciphertext[16+12:]
	var key []byte
	var err error

	key, err = gimmeKey(password, salt)
	if err != nil {
		return EasyCipher{}, err
	}

	var ec EasyCipher = EasyCipher{
		Password:   password,
		Key:        key,
		Salt:       salt,
		Iv:         iv,
		Ciphertext: cleanciphertext,
		Plaintext:  nil,
	}

	return ec, nil
}

// Encrypt() wraps a call to gcm.Seal(), and encrypts the plaintext. By design
// this will only work once - calling EasyCipher{}.Encrypt() when ec.Ciphertext
// is no longer nil will cause the program to panic(). I did this out of caution
// because I'm not sure that allowing anyone to reuse EasyCiphers{} won't
// cause problems or introduce security vulnerabilities.
func (ec *EasyCipher) Encrypt() {
	var gcm cipher.AEAD
	var ciphertext []byte
	var err error

	gcm, err = gimmeGCMCipher(ec.Key)
	if err != nil {
		panic(err)
	}

	if ec.Ciphertext != nil {
		err = errors.New("ec.Ciphertext not nil! Create a new EC instead of reusing this one")
		panic(err)
	}

	ciphertext = gcm.Seal(nil, ec.Iv, ec.Plaintext, nil)
	var siv []byte = append(ec.Salt, ec.Iv...)
	var sivcipher []byte = append(siv, ciphertext...)
	ec.Ciphertext = sivcipher
}

// Decrypt() wraps a call to gcm.Open(), and decrypts the ciphertext. By design
// this will only work once - calling EasyCipher{}.Decrypt() when ec.Plaintext
// is no longer nil will cause the program to panic(). I did this out of caution
// because I'm not sure that allowing anyone to reuse EasyCiphers{} won't
// cause problems or introduce security vulnerabilities.
func (ec *EasyCipher) Decrypt() {
	var gcm cipher.AEAD
	var dirtycipher []byte
	var err error

	gcm, err = gimmeGCMCipher(ec.Key)
	if err != nil {
		panic(err)
	}

	if ec.Plaintext != nil {
		err = errors.New("ec.Plaintext not nil! Create a new EC instead of reusing this one")
		panic(err)
	}

	dirtycipher, err = gcm.Open(nil, ec.Iv, ec.Ciphertext, nil)
	if err != nil {
		panic(err)
	}
	ec.Plaintext = dirtycipher
}

func gimmeGCMCipher(key []byte) (cipher.AEAD, error) {
	var block cipher.Block
	var aesgcm cipher.AEAD
	var err error

	block, err = aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	aesgcm, err = cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesgcm, nil
}

func gimmeSalt() ([]byte, error) {
	var err error
	var salt []byte = make([]byte, aes.BlockSize)

	_, err = io.ReadFull(rand.Reader, salt)

	if err != nil {
		return nil, err
	}

	return salt, nil
}

func gimmeIV() ([]byte, error) {
	// IV needs to be 12 bytes otherwise the GCM will complain, I have no
	// idea why it would I'm just following directions
	var err error
	var iv []byte = make([]byte, 12)
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return nil, err
	}

	return iv, nil
}

func gimmeKey(password string, salt []byte) ([]byte, error) {
	var key []byte
	var err error

	// for pbkdf2.Key() to return the same key given a password, the salt
	// used needs to be the same.
	key, err = pbkdf2.Key(sha256.New, password, salt, 4096, 32)
	if err != nil {
		return nil, err
	} else {
		return key, err
	}
}
