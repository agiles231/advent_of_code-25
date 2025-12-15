package tachyon

import (
	"fmt"
)

type TachyonSimulation struct {
	simulationStep int
	splittersHit   int
	manifold       TachyonManifold
	tachyonBeams   []*TachyonBeam
	grid           Grid
}

func (ts *TachyonSimulation) Tick() {
	updatedTachyonBeams := make([]*TachyonBeam, 0)
	numBeams := len(ts.tachyonBeams)
	fmt.Printf("Step %2v\n", ts.simulationStep)
	for i := 0; i < numBeams; i++ {
		tachyonBeam := ts.tachyonBeams[i]
		tachyonBeam.Tick()
		// fmt.Printf("Beam %2v  ", i)
		// tachyonBeam.Print()
		ii, jj, _, _ := tachyonBeam.GetPositionAndDirection()
		splitter, ok := ts.manifold.splitterLocations[ii][jj]
		if ok {
			splitter.isHit = true
			beams := ts.splitBeam(i)
			updatedTachyonBeams = append(updatedTachyonBeams, beams...)
		} else {
			updatedTachyonBeams = append(updatedTachyonBeams, tachyonBeam)
		}
	}
	ts.tachyonBeams = updatedTachyonBeams
	ts.dedupeBeams()
	ts.splittersHit = ts.CountSplits()
	ts.simulationStep++
}

func (ts *TachyonSimulation) CompleteSimulation() {
	for ts.simulationStep < ts.manifold.height {
		ts.Tick()
	}
}

func (ts TachyonSimulation) Print() {
	fmt.Printf("Simulation Step %5v, ", ts.simulationStep)
	fmt.Printf("Splitters Hit   %5v, ", ts.splittersHit)
	fmt.Printf("Num Beams       %5v\n", len(ts.tachyonBeams))
	ts.DrawGrid()
}

func (ts *TachyonSimulation) CountSplits() int {
	count := 0
	for _, iMap := range ts.manifold.splitterLocations {
		for _, splitter := range iMap {
			if splitter.isHit {
				count++
			}
		}
	}
	return count
}

func (ts *TachyonSimulation) GetBeamCount() int {
	count := 0
	for i := range ts.tachyonBeams {
		beam := ts.tachyonBeams[i]
		count += beam.count
	}
	return count
}

func (ts *TachyonSimulation) PrintBeams() {
	for i := range ts.tachyonBeams {
		b := ts.tachyonBeams[i]
		b.Print()
	}
}

func (ts *TachyonSimulation) splitBeam(i int) []*TachyonBeam {
	b := ts.tachyonBeams[i]
	leftBeam := b.Duplicate()
	rightBeam := b.Duplicate()
	leftBeam.Update(leftBeam.locationI, leftBeam.locationJ-1, leftBeam.directionI, leftBeam.directionJ)
	rightBeam.Update(rightBeam.locationI, rightBeam.locationJ+1, rightBeam.directionI, rightBeam.directionJ)
	leftValid := ts.validateBeam(leftBeam)
	rightValid := ts.validateBeam(rightBeam)
	beams := make([]*TachyonBeam, 0)
	if leftValid {
		beams = append(beams, leftBeam)
	}
	if rightValid {
		beams = append(beams, rightBeam)
	}
	return beams
}

func (ts *TachyonSimulation) dedupeBeams() {
	tachyonBeamTracker := make(map[int]map[int]*TachyonBeam)
	for i := 0; i < len(ts.tachyonBeams); i++ {
		beam := ts.tachyonBeams[i]
		ii, jj, _, _ := beam.GetPositionAndDirection()
		beamMap, ok := tachyonBeamTracker[ii]
		if ok {
			existingBeam, ok2 := beamMap[jj]
			if ok2 {
				existingBeam.count += beam.count
			} else {
				beamMap[jj] = beam
			}
		} else {
			beamMap = make(map[int]*TachyonBeam)
			tachyonBeamTracker[ii] = beamMap
			beamMap[jj] = beam
		}
	}
	updatedBeams := make([]*TachyonBeam, 0)
	for _, beamMap := range tachyonBeamTracker {
		for _, beam := range beamMap {
			updatedBeams = append(updatedBeams, beam)
		}
	}
	ts.tachyonBeams = updatedBeams
}

func (ts *TachyonSimulation) validateBeam(beam *TachyonBeam) bool {
	i := beam.locationI
	j := beam.locationJ
	if i < 0 || i >= ts.manifold.height {
		return false
	}
	if j < 0 || j >= ts.manifold.width {
		return false
	}
	return true
}

func (ts *TachyonSimulation) DrawGrid() {
	height := ts.manifold.height
	width := ts.manifold.width
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if i == 0 && j == ts.manifold.startIndex {
				fmt.Print("S")
			} else if ts.SplitterHitAtIJ(i, j) {
				fmt.Print("%")
			} else if ts.HasTachyonSplitterAtIJ(i, j) {
				fmt.Print("^")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func (ts *TachyonSimulation) HasTachyonSplitterAtIJ(i, j int) bool {
	iMap, ok := ts.manifold.splitterLocations[i]
	if ok {
		_, ok2 := iMap[j]
		return ok2
	}
	return false
}

func (ts *TachyonSimulation) SplitterHitAtIJ(i, j int) bool {
	iMap, ok := ts.manifold.splitterLocations[i]
	if ok {
		splitter, ok2 := iMap[j]
		if ok2 {
			return splitter.isHit
		}
	}
	return false
}

func (ts *TachyonSimulation) HasBeamAtIJ(i, j int) bool {
	for i := 0; i < len(ts.tachyonBeams); i++ {
		beam := ts.tachyonBeams[i]
		if beam.locationI == i && beam.locationJ == j {
			return true
		}
	}
	return false
}
