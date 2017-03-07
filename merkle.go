package pool

import (
	"github.com/dist-ribut-us/crypto"
	"github.com/dist-ribut-us/merkle"
	"io/ioutil"
)

const (
	Dir      = "./poolData/"
	SaltFile = "salt.bin"
	saltLen  = 16
)

var ErrBadSetup = defineErr("Bad setup")

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
