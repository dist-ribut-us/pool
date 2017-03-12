package pool

import (
	"github.com/dist-ribut-us/crypto"
	"github.com/dist-ribut-us/errors"
	"github.com/dist-ribut-us/merkle"
	"io/ioutil"
)

const (
	// Dir is the directory for the pool data, including the Merkle forrest
	Dir = "./poolData/"
	// SaltFile is the name of the file in Dir where the salt is stored
	SaltFile = "salt.bin"
	saltLen  = 16
)

// ErrBadSetup will be returned if the merkle forrest ccannot be read
const ErrBadSetup = errors.String("Bad setup")

func openMerkle(passphrase []byte) (*merkle.Forest, error) {
	salt, err := ioutil.ReadFile(Dir + SaltFile)
	if err != nil {
		return nil, err
	}
	if len(salt) != saltLen {
		return nil, ErrBadSetup
	}
	key := crypto.Hash(passphrase, salt).Digest().Shared()
	return merkle.Open(Dir, key)
}
