package main

import (
	"fmt"
	//"github.com/dist-ribut-us/ipc"
	"github.com/dist-ribut-us/pool"
	"github.com/howeyc/gopass"
	"time"
)

func main() {
	passphrase := getPassphrase()
	p, err := pool.Open(passphrase)
	check(err)
	fmt.Println(p.List())
	p.Start()
	for {
		time.Sleep(time.Hour)
	}
}

func getPassphrase() []byte {
	fmt.Print("Password: ")
	pass, err := gopass.GetPasswd()
	check(err)
	return pass
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
