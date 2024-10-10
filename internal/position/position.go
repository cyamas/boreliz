package position

import "github.com/cyamas/boreliz/internal/hold"

type Position struct {
	leftHand  *hold.Hold
	rightHand *hold.Hold
	leftFoot  *hold.Hold
	rightFoot *hold.Hold
}

func New(lh, rh, lf, rf *hold.Hold) *Position {
	return &Position{
		leftHand:  lh,
		rightHand: rh,
		leftFoot:  lf,
		rightFoot: rf,
	}
}
