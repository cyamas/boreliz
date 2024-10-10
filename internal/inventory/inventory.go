package inventory

import (
	"errors"

	"github.com/cyamas/boreliz/internal/hold"
)

type Inventory struct {
	holds map[int]*hold.Hold
}

func New() *Inventory {
	return &Inventory{holds: make(map[int]*hold.Hold)}
}

func Load() *Inventory {
	inv := New()
	return inv
}

func (i *Inventory) AllHolds() []*hold.Hold {
	holds := []*hold.Hold{}
	for _, hold := range i.holds {
		holds = append(holds, hold)
	}
	return holds
}

func (i *Inventory) AddHold(hold *hold.Hold) {
	i.holds[hold.ID()] = hold
}

func (i *Inventory) GetHoldByID(id int) (*hold.Hold, error) {
	if hold, ok := i.holds[id]; ok {
		return hold, nil
	}
	return nil, errors.New("hold does not exist")
}
