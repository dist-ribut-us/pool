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
		p.Router.NetSenderPort = p.overlay
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
	p.Router.Register(p)
	go p.startAll()
	p.Router.Run()
}

// ServiceID for Pool service
func (*Pool) ServiceID() uint32 {
	return message.PoolService
}

// AddBeacon to help connect to network
func (p *Pool) AddBeacon(addr *rnet.Addr, key *crypto.XchgPub) {
	p.Router.
		Command(message.AddBeacon, key.Slice()).
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
			Command(message.Die, nil).
			To(prg.Port()).
			Send(nil)
	}
}

// QueryHandler for ipc queries to Pool service
func (p *Pool) QueryHandler(q ipcrouter.Query) {
	log.Info(log.Lbl("handling_query"))
	switch t := q.GetType(); t {
	case message.GetPort:
		name := q.BodyString()
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
		Send(func(r ipcrouter.Response) {
			ch <- crypto.XchgPubFromSlice(r.GetBody()).String()
		})
	return <-ch
}

// GetOverlayNetPort gets the network port
func (p *Pool) GetOverlayNetPort() uint32 {
	ch := make(chan uint32)
	p.Router.
		Query(message.GetPort, nil).
		To(p.overlay).
		Send(func(r ipcrouter.Response) {
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
		SetService(19860714).
		SendToNet(from, func(r ipcrouter.NetResponse) {
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
		Command(message.RandomKey, nil).
		To(p.overlay).
		Send(nil)
}

// OverlayStaticKey tells the overlay service to use a static key
func (p *Pool) OverlayStaticKey() {
	p.Router.
		Command(message.StaticKey, nil).
		To(p.overlay).
		Send(nil)
}
