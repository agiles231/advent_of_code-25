package tachyon

type TachyonManifold struct {
	startIndex        int
	splitterLocations map[int]map[int]*TachyonSplitter
	height            int
	width             int
}

type TachyonSplitter struct {
	isHit bool
}

func NewTachyonManifold(g Grid) (*TachyonManifold, error) {
	startIndex := 0
	height := g.Height()
	width := g.Width()
	splitterLocations := make(map[int]map[int]*TachyonSplitter)
	g.IterateGrid(func(i, j int, value rune) {
		if value == 'S' {
			startIndex = j
		}
		if value == '^' {
			_, ok := splitterLocations[i]
			if !ok {
				splitterLocations[i] = make(map[int]*TachyonSplitter)
			}
			splitterLocations[i][j] = &TachyonSplitter{isHit: false}
		}
	})
	return &TachyonManifold{startIndex: startIndex, splitterLocations: splitterLocations, height: height, width: width}, nil
}

func (t *TachyonManifold) NewSimulation() *TachyonSimulation {
	grid := make(Grid, t.height)
	for i := 0; i < t.height; i++ {
		grid[i] = make([]rune, t.width)
		for j := 0; j < t.width; j++ {
			_, splitterExists := t.splitterLocations[i][j]
			if j == 0 && i == t.startIndex {
				grid[i][j] = 'S'
			} else if splitterExists {
				grid[i][j] = '^'
			} else {
				grid[i][j] = '.'
			}
		}
	}
	startIndex := t.startIndex
	startBeam := NewTachyonBeam(0, startIndex, 1, 0, 1)
	beams := make([]*TachyonBeam, 1)
	beams[0] = startBeam
	return &TachyonSimulation{
		simulationStep: 0,
		manifold:       *t,
		tachyonBeams:   beams,
		grid:           grid,
	}
}
