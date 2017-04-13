package pool

import (
	"github.com/dist-ribut-us/crypto"
	"github.com/dist-ribut-us/ipcrouter"
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
	Router   *ipcrouter.Router
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
	p.Router.Register(message.PoolService, p.handler)
	go p.startAll()
	p.Router.Run()
}

// AddBeacon to help connect to network
func (p *Pool) AddBeacon(addr *rnet.Addr, key *crypto.XchgPub) {
	p.Router.
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
		p.Router.
			Base(message.Die, nil).
			To(prg.Port()).
			Send(nil)
	}
}

func (p *Pool) handler(b *ipcrouter.Base) {
	if b.IsQuery() {
		go p.handleQuery(b)
	} else {
		log.Info(log.Lbl("pool_unknown_type"), b.GetType())
	}
}

// handleQuery takes a wrapper and responds to it's query
func (p *Pool) handleQuery(q *ipcrouter.Base) {
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
	p.Router.
		Query(message.GetPubKey, nil).
		To(p.overlay).
		Send(func(r *ipcrouter.Base) {
			ch <- crypto.XchgPubFromSlice(r.Body).String()
		})
	return <-ch
}

// GetOverlayNetPort gets the network port
func (p *Pool) GetOverlayNetPort() uint32 {
	ch := make(chan uint32)
	p.Router.
		Query(message.GetPort, nil).
		To(p.overlay).
		Send(func(r *ipcrouter.Base) {
			ch <- r.BodyToUint32()
		})
	return <-ch
}

// GetIP is a temporary method for testing.
func (p *Pool) GetIP(from *rnet.Addr) *rnet.Addr {
	ch := make(chan *rnet.Addr)
	log.Info("sending_get_ip_request")
	p.Router.
		Query(message.GetIP, nil).
		ToNet(p.overlay, from, 19860714).
		Send(func(r *ipcrouter.Base) {
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

// OverlayRandomKey tells the overlay service to use a random key
func (p *Pool) OverlayRandomKey() {
	p.Router.
		Base(message.RandomKey, nil).
		To(p.overlay).
		Send(nil)
}

// OverlayStaticKey tells the overlay service to use a static key
func (p *Pool) OverlayStaticKey() {
	p.Router.
		Base(message.StaticKey, nil).
		To(p.overlay).
		Send(nil)
}
