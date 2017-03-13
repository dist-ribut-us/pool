package pool

import (
	"github.com/dist-ribut-us/crypto"
	"github.com/dist-ribut-us/errors"
	"github.com/dist-ribut-us/merkle"
	"github.com/dist-ribut-us/prog"
	"io/ioutil"
)

const (
	saltFile = "salt.bin"
	saltLen  = 16
)

// LogFile for pool log
var LogFile = prog.Root() + "pool.log"
var dir = prog.Root() + "poolData/"

// Dir returns the directory used for pool data
func Dir() string { return dir }

// SaltFile gets the location of the salt file
func SaltFile() string { return saltFile }

// ErrBadSetup will be returned if the merkle forrest ccannot be read
const ErrBadSetup = errors.String("Bad setup")

func openMerkle(passphrase []byte) (*merkle.Forest, error) {
	salt, err := getSalt()
	if err != nil {
		return nil, err
	}
	key := crypto.Hash(passphrase, salt).Digest().Shared()
	return merkle.Open(dir, key)
}

// IsSetup checks if the salt file exists as an litmus for if it's setup.
func IsSetup() bool {
	_, err := getSalt()
	return err == nil
}

func getSalt() ([]byte, error) {
	salt, err := ioutil.ReadFile(dir + saltFile)
	if err != nil {
		return nil, err
	}
	if len(salt) != saltLen {
		return nil, ErrBadSetup
	}
	return salt, nil
}
