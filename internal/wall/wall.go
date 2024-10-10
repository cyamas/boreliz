package wall

import "fmt"

type Wall struct {
	ID         string
	Angle      int
	Rows, Cols int
	Grid       []*Bolt
}

type Bolt struct {
	ID string
}

func New() *Wall {
	return &Wall{}
}

func (w *Wall) CreateGrid() {
	for i := range w.Rows {
		for j := range w.Cols {
			id := fmt.Sprintf("%d-%d", i, j)
			bolt := &Bolt{ID: id}
			w.Grid = append(w.Grid, bolt)
		}
	}
}
