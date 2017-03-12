package main

import (
	"fmt"
	"github.com/dist-ribut-us/crypto"
	"github.com/dist-ribut-us/log"
	"github.com/dist-ribut-us/pool"
	"github.com/howeyc/gopass"
)

func main() {
	log.Contents = log.Truncate
	log.Panic(log.ToFile())
	log.Go()

	log.Info(log.Lbl("starting_pool"))

	var p *pool.Pool
	var err error
	for {
		passphrase := getPassphrase()
		p, err = pool.Open(passphrase)
		if err == nil {
			break
		} else if err == crypto.ErrDecryptionFailed {
			continue
		} else {
			log.Panic(err)
		}
	}
	p.Start()
	log.Info("pool_listening")
	for msg := range p.Chan() {
		b, err := msg.ToBase()
		if log.Error(err) {
			continue
		}
		if b.IsQuery() {
			go p.HandleQuery(b)
		} else {
			log.Info(log.Lbl("pool_unknown_type"), b.Type)
		}
	}
}

func getPassphrase() []byte {
	fmt.Print("Password: ")
	pass, err := gopass.GetPasswd()
	log.Panic(err)
	return pass
}
