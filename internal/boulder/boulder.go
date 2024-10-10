package boulder

import (
	"github.com/cyamas/boreliz/internal/hold"
	"github.com/cyamas/boreliz/internal/move"
)

type Boulder struct {
	id    int
	name  string
	holds map[*hold.Hold]bool
	grade int
	moves []*move.Move
}

func New() *Boulder {
	return &Boulder{}
}

func (b *Boulder) Set() {

}

func (b *Boulder) AddHold(hold *hold.Hold) {
	b.holds[hold] = true
}

func (b *Boulder) AddMove() {
	newMove := move.New()
	b.moves = append(b.moves, newMove)
}

func (b *Boulder) ID() int {
	return b.id
}

func (b *Boulder) SetName(name string) {
	b.name = name
}

func (b *Boulder) Name() string {
	return b.name
}

func (b *Boulder) SetGrade(grade int) {
	b.grade = grade
}

func (b *Boulder) Grade() int {
	return b.grade
}
