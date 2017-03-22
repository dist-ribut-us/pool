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

func userSetup() *pool.Pool {
	passphrase := setPassphrase()
	return setup(passphrase)
}

func beaconSetup() *pool.Pool {
	p := setup(nil)
	addBeacon(p)
	return p
}

func setup(passphrase []byte) *pool.Pool {
	buildMerkle(passphrase)
	p, err := pool.Open(passphrase)
	log.Panic(err)
	addOverlay(p)
	addDHT(p)
	return p
}

func addOverlay(p *pool.Pool) {
	p.Add(&pool.Program{
		Name:     "Overlay",
		Location: "ribut.overlay",
		UI:       false,
		Key:      crypto.RandomShared().Slice(),
		Port32:   uint32(rnet.RandomPort()),
		Start:    true,
	})
}

func addDHT(p *pool.Pool) {
	p.Add(&pool.Program{
		Name:     "DHT",
		Location: "ribut.dht",
		UI:       false,
		Key:      crypto.RandomShared().Slice(),
		Port32:   uint32(rnet.RandomPort()),
		Start:    true,
	})
}

func addBeacon(p *pool.Pool) {
	p.Add(&pool.Program{
		Name:     "Beacon",
		Location: "ribut.beacon",
		UI:       false,
		Key:      crypto.RandomShared().Slice(),
		Port32:   uint32(rnet.RandomPort()),
		Start:    true,
	})
}

func setPassphrase() []byte {
	for {
		fmt.Print("\nPassword: ")
		pass1, err := gopass.GetPasswd()
		log.Panic(err)
		fmt.Print("\nAgain: ")
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
	dir := pool.Dir()
	log.Panic(os.MkdirAll(dir, 0777))
	// now that directory exists, we can set logging up
	log.Panic(log.ToFile(pool.LogFile))
	log.Info(log.Lbl("running_setup"))
	saltFile, err := os.Create(dir + pool.SaltFile())
	log.Panic(err)
	defer func() { log.Panic(saltFile.Close()) }()
	_, err = saltFile.Write(salt)
	log.Panic(err)
	key := crypto.Hash(passphrase, salt).Digest().Shared()
	forest, err := merkle.Open(dir, key)
	log.Panic(err)
	forest.Close()
}
