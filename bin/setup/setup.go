package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/dist-ribut-us/crypto"
	"github.com/dist-ribut-us/log"
	"github.com/dist-ribut-us/merkle"
	"github.com/dist-ribut-us/pool"
	"github.com/dist-ribut-us/rnet"
	"github.com/howeyc/gopass"
	"os"
)

func main() {
	log.Panic(log.ToFile())
	log.Go()

	passphrase := setPassphrase()

	buildMerkle(passphrase)
	p, err := pool.Open(passphrase)
	log.Panic(err)

	p.Add(&pool.Program{
		Name:     "Overlay",
		Location: "./overlay",
		UI:       false,
		Key:      crypto.RandomShared().Slice(),
		Port32:   uint32(rnet.RandomPort()),
		Start:    true,
	})
	p.Add(&pool.Program{
		Name:     "DHT",
		Location: "./dht",
		UI:       false,
		Key:      crypto.RandomShared().Slice(),
		Port32:   uint32(rnet.RandomPort()),
		Start:    true,
	})

	fmt.Println("Pool is setup")
}

func setPassphrase() []byte {
	for {
		fmt.Println("Password: ")
		pass1, err := gopass.GetPasswd()
		log.Panic(err)
		fmt.Println("Again: ")
		pass2, err := gopass.GetPasswd()
		log.Panic(err)
		if bytes.Equal(pass1, pass2) {
			return pass1
		}
		fmt.Println("Passwords did not match")
	}
}

func buildMerkle(passphrase []byte) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	log.Panic(err)
	log.Panic(os.MkdirAll(pool.Dir, 0777))
	saltFile, err := os.Create(pool.Dir + pool.SaltFile)
	log.Panic(err)
	defer func() { log.Panic(saltFile.Close()) }()
	_, err = saltFile.Write(salt)
	log.Panic(err)
	key := crypto.Hash(passphrase, salt).Digest().Shared()
	forest, err := merkle.Open(pool.Dir, key)
	log.Panic(err)
	forest.Close()
}
