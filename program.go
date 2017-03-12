package pool

import (
	"github.com/dist-ribut-us/rnet"
	"strconv"
)

// Port gets the program port
func (p *Program) Port() rnet.Port { return rnet.Port(p.Port32) }

// PortStr gets the port as a string (without :)
func (p *Program) PortStr() string { return strconv.Itoa(int(p.Port32)) }
