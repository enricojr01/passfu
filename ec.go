package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

type EasyCipher struct {
	Password   string
	Key        []byte
	Salt       []byte
	Iv         []byte
	Ciphertext []byte
	Plaintext  []byte
}

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

func NewFromCiphertext(ciphertext []byte, password string) (EasyCipher, error) {
	var salt []byte = ciphertext[:16]
	var iv []byte = ciphertext[16 : 16+12]
	var cleanciphertext []byte = ciphertext[16+12:]
	var key []byte
	var err error

	key, err = gimmeKey(password, salt)
	if err != nil {
		return EasyCipher{}, err
	}

	// Ciphertext needs to include the salt + iv otherwise authentication
	// will fail.
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

func (ec *EasyCipher) Encrypt() {
	// implementation notes:
	// take plaintext, encrypt first, THEN append the salt + iv
	// encrypting the salt + iv will cause message auth to fail
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
