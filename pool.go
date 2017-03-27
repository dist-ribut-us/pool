package pool

import (
	"github.com/dist-ribut-us/crypto"
	"github.com/dist-ribut-us/ipc"
	"github.com/dist-ribut-us/log"
	"github.com/dist-ribut-us/merkle"
	"github.com/dist-ribut-us/message"
	"github.com/dist-ribut-us/rnet"
	"github.com/dist-ribut-us/serial"
	"github.com/golang/protobuf/proto"
	"os/exec"
	"time"
)

// Port that pool will run on, should be moved to config
var Port = rnet.Port(3000)

var progBkt = []byte("programs")

// Pool coordinates all local resources for dist.ribut.us
type Pool struct {
	forrest  *merkle.Forest
	programs map[string]*Program
	ipc      *ipc.Proc
	overlay  rnet.Port
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

	if overlay, ok := p.programs["Overlay"]; ok {
		p.overlay = overlay.Port()
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

// Run the ipc listener and all the progams designated to start
func (p *Pool) Run() {
	var err error
	p.ipc, err = ipc.New(Port)
	if log.Error(err) {
		return
	}
	p.ipc.Handler(p.handler)
	go p.startAll()
	p.ipc.Run()
}

// AddBeacon to help connect to network
func (p *Pool) AddBeacon(addr *rnet.Addr, key *crypto.XchgPub) {
	p.ipc.
		Base(message.AddBeacon, key.Slice()).
		SetAddr(addr).
		To(p.overlay).
		Send(nil)
}

func (p *Pool) startAll() {
	p.ExitAll()
	time.Sleep(time.Millisecond)
	for _, prg := range p.programs {
		if !prg.Start {
			continue
		}
		go p.run(prg)
	}
}

func (p *Pool) run(prg *Program) {
	lg := log.Child(prg.Name)
	log.Info(log.Lbl("starting"), prg.GetLocation())
	cmd := exec.Command(prg.GetLocation(), prg.PortStr(), Port.RawStr(), crypto.SymmetricFromSlice(prg.Key).String())
	out, err := cmd.CombinedOutput()
	if lg.Error(err) {
		lg.Info(string(out))
	}
}

// ExitAll sends an exit message to all service, whether they are running or
// not (closes stray programs from previous)
func (p *Pool) ExitAll() {
	for _, prg := range p.programs {
		p.ipc.
			Base(message.Die, nil).
			To(prg.Port()).
			Send(nil)
	}
}

func (p *Pool) handler(b *ipc.Base) {
	if b.IsQuery() {
		go p.handleQuery(b)
	} else {
		log.Info(log.Lbl("pool_unknown_type"), b.GetType())
	}
}

// handleQuery takes a wrapper and responds to it's query
func (p *Pool) handleQuery(q *ipc.Base) {
	log.Info(log.Lbl("handling_query"))
	switch t := q.GetType(); t {
	case message.GetPort:
		name := string(q.Body)
		prg, ok := p.programs[name]
		if !ok {
			log.Info(log.Lbl("bad_port_request"), name, q.Port())
			return
		}
		q.Respond(serial.MarshalUint32(prg.Port32, nil))
	default:
		log.Info("unknown_query", t)
	}
}

// GetOverlayPubKey gets the public key for the network
func (p *Pool) GetOverlayPubKey() string {
	ch := make(chan string)
	p.ipc.
		Query(message.GetPubKey, nil).
		To(p.overlay).
		Send(func(r *ipc.Base) {
			ch <- crypto.XchgPubFromSlice(r.Body).String()
		})
	return <-ch
}

// GetOverlayNetPort gets the network port
func (p *Pool) GetOverlayNetPort() uint32 {
	ch := make(chan uint32)
	p.ipc.
		Query(message.GetPort, nil).
		To(p.overlay).
		Send(func(r *ipc.Base) {
			ch <- r.BodyToUint32()
		})
	return <-ch
}

// GetIP is a temporary method for testing.
func (p *Pool) GetIP(from *rnet.Addr) *rnet.Addr {
	ch := make(chan *rnet.Addr)
	log.Info("sending_get_ip_request")
	p.ipc.
		Query(message.GetIP, nil).
		ToNet(p.overlay, from, 19860714).
		Send(func(r *ipc.Base) {
			var addrpb message.Addrpb
			err := r.Unmarshal(&addrpb)
			if log.Error(err) {
				return
			}
			ch <- addrpb.GetAddr()
		})

	var r *rnet.Addr
	select {
	case addr := <-ch:
		r = addr
	case <-time.After(time.Millisecond * 100):
		r = nil
	}
	return r
}
