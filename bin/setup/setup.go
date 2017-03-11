package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/dist-ribut-us/crypto"
	"github.com/dist-ribut-us/merkle"
	"github.com/dist-ribut-us/pool"
	"github.com/dist-ribut-us/rnet"
	"github.com/howeyc/gopass"
	"os"
)

func main() {
	passphrase := setPassphrase()

	f, err := buildMerkle(passphrase)
	check(err)
	f.Close()

	p, err := pool.Open(passphrase)

	key := crypto.RandomShared()

	overlay := &pool.Program{
		Name:     "Overlay",
		Location: "./overlay",
		UI:       false,
		Key:      key.Slice(),
		Port32:   uint32(rnet.RandomPort()),
		Start:    true,
	}
	p.Add(overlay)

	fmt.Println("Pool is setup")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func setPassphrase() []byte {
	for {
		fmt.Println("Password: ")
		pass1, err := gopass.GetPasswd()
		if err != nil {
			panic(err)
		}
		fmt.Println("Again: ")
		pass2, err := gopass.GetPasswd()
		if err != nil {
			panic(err)
		}
		if bytes.Equal(pass1, pass2) {
			return pass1
		}
		fmt.Println("Passwords did not match")
	}
}

func buildMerkle(passphrase []byte) (*merkle.Forest, error) {
	salt := make([]byte, 16)
	rand.Read(salt)
	os.MkdirAll(pool.Dir, 0777)
	saltFile, err := os.Create(pool.Dir + pool.SaltFile)
	if err != nil {
		return nil, err
	}
	defer saltFile.Close()
	saltFile.Write(salt)
	key := crypto.Hash(passphrase, salt).Digest().Shared()
	return merkle.Open(pool.Dir, key)
}
