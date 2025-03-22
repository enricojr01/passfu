package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

type EasyCipher struct {
	Password   string
	Key        []byte
	Salt       []byte
	Iv         []byte
	Ciphertext []byte
}

func (ec *EasyCipher) Encrypt(plaintext string) {
	var gcm cipher.AEAD
	var plaintextByte []byte = []byte(plaintext)
	var err error

	gcm, err = gimmeGCMCipher(ec.Key)
	if err != nil {
		panic(err)
	}

	ec.Ciphertext = gcm.Seal(nil, ec.Iv, plaintextByte, nil)
}

func (ec *EasyCipher) ExportCiphertext() []byte {
	var saltiv []byte = append(ec.Salt, ec.Iv...)
	var saltivcipher []byte = append(saltiv, ec.Ciphertext...)
	return saltivcipher
}

func (ec *EasyCipher) Decrypt() []byte {
	var gcm cipher.AEAD
	var plaintextByte []byte
	var err error

	gcm, err = gimmeGCMCipher(ec.Key)
	if err != nil {
		panic(err)
	}

	plaintextByte, err = gcm.Open(nil, ec.Iv, ec.Ciphertext, nil)
	if err != nil {
		panic(err)
	}

	return plaintextByte
}

func NewECFromCiphertext(ciphertext []byte, password string) EasyCipher {
	var salt []byte = ciphertext[:16]
	var iv []byte = ciphertext[16 : 16+12]
	var key []byte
	var err error

	key, err = gimmeKey(password, salt)
	if err != nil {
		panic(err)
	}

	var ec EasyCipher = EasyCipher{
		Password:   password,
		Key:        key,
		Salt:       salt,
		Iv:         iv,
		Ciphertext: ciphertext[16+12:],
	}

	return ec
}

func NewEC(password string) (EasyCipher, error) {
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
	}

	return ec, nil
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

	key, err = pbkdf2.Key(sha256.New, password, salt, 4096, 32)
	if err != nil {
		return nil, err
	} else {
		return key, err
	}
}
