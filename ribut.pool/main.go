package main

import (
	"fmt"
	"github.com/dist-ribut-us/crypto"
	"github.com/dist-ribut-us/log"
	"github.com/dist-ribut-us/pool"
	"github.com/dist-ribut-us/prog"
	"github.com/dist-ribut-us/rnet"
	"github.com/howeyc/gopass"
	"github.com/urfave/cli"
	"os"
	"time"
)

func main() {
	log.Contents = log.Truncate
	log.Go()

	app := cli.NewApp()
	app.Name = "ribut.pool"
	app.Usage = "Run dist.ribut.us"
	app.Action = func(c *cli.Context) error {
		var p *pool.Pool
		if !pool.IsSetup() {
			p = userSetup()
		} else {
			log.Panic(log.ToFile(pool.LogFile))
			p = userOpen()
		}

		log.Info(log.Lbl("starting_pool"))

		go p.Run()
		runCLI(p)
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:  "beacon",
			Usage: "Run dist.ribut.us as a beacon",
			Action: func(c *cli.Context) error {
				var p *pool.Pool
				if !pool.IsSetup() {
					p = beaconSetup()
				} else {
					log.Panic(log.ToFile(pool.LogFile))
					p = beaconOpen()
				}

				log.Info(log.Lbl("starting_pool"))

				go func() {
					time.Sleep(time.Millisecond * 50)
					fmt.Println(p.GetOverlayNetPort(), p.GetOverlayPubKey())
				}()

				p.Run()
				return nil
			},
		},
	}

	app.Run(os.Args)

}

func getPassphrase() []byte {
	fmt.Print("Password: ")
	pass, err := gopass.GetPasswd()
	log.Panic(err)
	return pass
}

func userOpen() *pool.Pool {
	for {
		passphrase := getPassphrase()
		p, err := pool.Open(passphrase)
		if err == nil {
			return p
		} else if err == crypto.ErrDecryptionFailed {
			continue
		} else {
			log.Panic(err)
		}
	}
}

func beaconOpen() *pool.Pool {
	p, err := pool.Open(nil)
	log.Panic(err)
	return p
}

type command struct {
	desc    string
	handler func([]string, *pool.Pool)
}

var commands = map[string]command{
	"exit": {
		desc: "Exits pool",
		handler: func(s []string, p *pool.Pool) {
			p.ExitAll()
			time.Sleep(time.Millisecond)
			os.Exit(0)
		},
	},
	"beacon": {
		desc: "Add Beacon: url:port key",
		handler: func(s []string, p *pool.Pool) {
			if len(s) < 3 {
				fmt.Println("Require <url:port> <key>")
				return
			}
			addr, err := rnet.ResolveAddr(s[1])
			if err != nil {
				fmt.Println(err)
				return
			}
			pub, err := crypto.PubFromString(s[2])
			if err != nil {
				fmt.Println(err)
				return
			}
			p.AddBeacon(addr, pub)
		},
	},
}

func runCLI(p *pool.Pool) {
	// not in initilization to avoid initi loop
	commands["help"] = command{
		desc: "Print all commands",
		handler: func(s []string, p *pool.Pool) {
			for name, cmd := range commands {
				fmt.Println(name, ":", cmd.desc)
			}
		},
	}
	for input := prog.ReadStdin("pool> "); true; input = prog.ReadStdin("pool> ") {
		if len(input) < 1 {
			continue
		}
		if cmd, ok := commands[input[0]]; ok {
			cmd.handler(input, p)
		} else {
			fmt.Println(input[0], "is not a valid command")
			commands["help"].handler(input, p)
		}
	}
}
