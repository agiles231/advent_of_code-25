package tachyon

import (
	"fmt"
)

type TachyonBeam struct {
	locationI  int
	locationJ  int
	directionI int
	directionJ int
	count      int
}

func NewTachyonBeam(i, j, di, dj, count int) *TachyonBeam {
	return &TachyonBeam{
		locationI:  i,
		locationJ:  j,
		directionI: di,
		directionJ: dj,
		count:      count,
	}
}

func (b *TachyonBeam) Tick() {
	b.locationI += b.directionI
	b.locationJ += b.directionJ
}

func (b *TachyonBeam) Count() int {
	return b.count
}

func (b *TachyonBeam) Print() {
	fmt.Printf("Count: %8v, location: [%2v, %2v], direction[%2v, %2v]\n", b.count, b.locationI, b.locationJ, b.directionI, b.directionJ)
}

func (b *TachyonBeam) Duplicate() *TachyonBeam {
	return &TachyonBeam{
		locationI:  b.locationI,
		locationJ:  b.locationJ,
		directionI: b.directionI,
		directionJ: b.directionJ,
		count:      b.count,
	}
}

func (b *TachyonBeam) Update(x, y, dx, dy int) {
	b.locationI = x
	b.locationJ = y
	b.directionI = dx
	b.directionJ = dy
}

func (b *TachyonBeam) GetPositionAndDirection() (int, int, int, int) {
	return b.locationI, b.locationJ, b.directionI, b.directionJ
}
