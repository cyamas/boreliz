package climbs

import (
	"errors"

	"github.com/cyamas/boreliz/internal/boulder"
)

type List struct {
	boulders map[int]*boulder.Boulder
}

func (l *List) GetAllBoulders() map[int]*boulder.Boulder {
	return l.boulders
}

func (l *List) GetBoulderByID(id int) (*boulder.Boulder, error) {
	if boulder, ok := l.boulders[id]; ok {
		return boulder, nil
	}
	return nil, errors.New("boulder does not exist")
}
