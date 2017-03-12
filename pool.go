package pool

import (
	"fmt"
	"github.com/dist-ribut-us/crypto"
	"github.com/dist-ribut-us/ipc"
	"github.com/dist-ribut-us/log"
	"github.com/dist-ribut-us/merkle"
	"github.com/dist-ribut-us/rnet"
	"github.com/dist-ribut-us/serial"
	"github.com/golang/protobuf/proto"
	"os/exec"
)

// Port that pool will run on, should be moved to config
var Port = rnet.Port(3000)

var progBkt = []byte("programs")

// Pool coordinates all local resources for dist.ribut.us
type Pool struct {
	forrest  *merkle.Forest
	programs map[string]*Program
	ipc      *ipc.Proc
}

// Open the merkle forrest and setup the pool
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

// Add a program
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

// Start the ipc listener and all the progams designated to start
func (p *Pool) Start() {
	var err error
	p.ipc, err = ipc.RunNew(Port)
	log.Error(err)
	for _, prg := range p.programs {
		if !prg.Start {
			continue
		}
		go func(prg *Program) {
			lg := log.Child(prg.Name)
			fmt.Println(prg.GetLocation(), prg.PortStr(), Port.RawStr(), crypto.SharedFromSlice(prg.Key).String())
			cmd := exec.Command(prg.GetLocation(), prg.PortStr(), Port.RawStr(), crypto.SharedFromSlice(prg.Key).String())
			out, err := cmd.CombinedOutput()
			if lg.Error(err) {
				lg.Info(string(out))
			}
		}(prg)
	}
}

// Chan gets the ipc channel
func (p *Pool) Chan() <-chan *ipc.Message {
	return p.ipc.Chan()
}

// HandleQuery takes a wrapper and responds to it's query
func (p *Pool) HandleQuery(q *ipc.Base) {
	log.Info(log.Lbl("handling_query"))
	switch q.Type {
	case ipc.TPort:
		name := string(q.Body)
		prg, ok := p.programs[name]
		if !ok {
			log.Info(log.Lbl("bad_port_request"), name, q.Port())
			return
		}
		q.Respond(serial.MarshalUint32(prg.Port32, nil))
	}
}
