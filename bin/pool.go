package main

import (
	"fmt"
	"github.com/dist-ribut-us/crypto"
	"github.com/dist-ribut-us/log"
	"github.com/dist-ribut-us/pool"
	"github.com/howeyc/gopass"
	"time"
)

func main() {
	log.Panic(log.ToFile())
	log.Go()

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
	fmt.Println(p.List())
	p.Start()
	for {
		time.Sleep(time.Hour)
	}
}

func getPassphrase() []byte {
	fmt.Print("Password: ")
	pass, err := gopass.GetPasswd()
	log.Panic(err)
	return pass
}
