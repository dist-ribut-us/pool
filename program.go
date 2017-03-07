package pool

import (
	"fmt"
)

func (p *Program) PortStr() string { return fmt.Sprintf("%d", p.Port) }
