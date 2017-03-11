package pool

import (
	"github.com/dist-ribut-us/rnet"
	"strconv"
)

func (p *Program) Port() rnet.Port { return rnet.Port(p.Port32) }
func (p *Program) PortStr() string { return strconv.Itoa(int(p.Port32)) }
