package move

import (
	"github.com/cyamas/boreliz/internal/hold"
	"github.com/cyamas/boreliz/internal/position"
)

const (
	START = iota
	MID
	END
)

type Move struct {
	startPos *position.Position
	midPos   *position.Position
	endPos   *position.Position
}

func New() *Move {
	return &Move{}
}

func (m *Move) SetPosition(pos int, lh, rh, lf, rf *hold.Hold) {
	switch pos {
	case START:
		m.startPos = position.New(lh, rh, lf, rf)
	case MID:
		m.midPos = position.New(lh, rh, lf, rf)
	case END:
		m.endPos = position.New(lh, rh, lf, rf)
	}
}
