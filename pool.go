package pool

import (
	"fmt"
	"github.com/dist-ribut-us/crypto"
	"github.com/dist-ribut-us/merkle"
	"github.com/golang/protobuf/proto"
	"os/exec"
)

// this needs to be moved to config
const Port = "3000"

// this allows errors to be defined as const instead of var
type defineErr string

func (d defineErr) Error() string {
	return string(d)
}

var progBkt = []byte("programs")

type Pool struct {
	forrest  *merkle.Forest
	programs map[string]*Program
}

func Open(passphrase []byte) (*Pool, error) {
	f, err := openMerkle(passphrase)
	if err != nil {
		return nil, err
	}
	p := &Pool{
		forrest:  f,
		programs: make(map[string]*Program),
	}

	for k, v, err := f.First(progBkt); k != nil && err == nil; k, v, err = f.Next(progBkt, k) {
		var prg Program
		if err = proto.Unmarshal(v, &prg); err != nil {
			return nil, err
		}
		p.programs[prg.Name] = &prg
	}

	return p, nil
}

func (p *Pool) Len() int {
	return len(p.programs)
}

func (p *Pool) List() []string {
	var prgs []string
	for prg := range p.programs {
		prgs = append(prgs, prg)
	}
	return prgs
}

func (p *Pool) Status() string {
	return "OK"
}

func (p *Pool) Add(prog *Program) error {
	//todo: check if already in map
	data, err := proto.Marshal(prog)
	if err != nil {
		return err
	}
	p.forrest.SetValue(progBkt, []byte(prog.Name), data)
	p.programs[prog.Name] = prog
	return nil
}

func (p *Pool) Start() {
	for _, prg := range p.programs {
		if !prg.Start {
			continue
		}
		fmt.Println(prg.GetLocation(), prg.PortStr(), Port, crypto.SharedFromSlice(prg.Key).String())
		err := exec.Command(prg.GetLocation(), prg.PortStr(), Port, crypto.SharedFromSlice(prg.Key).String()).Run()
		if err != nil {
			fmt.Println("Start Error: ", err)
		}
	}
}
