package main

import (
	"fmt"
	"github.com/dist-ribut-us/crypto"
	"github.com/dist-ribut-us/ipc"
	"github.com/dist-ribut-us/log"
	"github.com/dist-ribut-us/pool"
	"github.com/howeyc/gopass"
)

func main() {
	log.Panic(log.ToFile())
	log.Go()

	log.Info("Starting Pool")

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
		w, err := msg.Unwrap()
		if log.Error(err) {
			continue
		}
		switch w.Type {
		case ipc.Type_QUERY:
			go p.HandleQuery(w)
		default:
			log.Info(log.Lbl("pool_unknown_type"), w.Type)
		}
	}
}

func getPassphrase() []byte {
	fmt.Print("Password: ")
	pass, err := gopass.GetPasswd()
	log.Panic(err)
	return pass
}
