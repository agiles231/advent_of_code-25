package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

/*
1. Convert list of points into list of lines and points
2. Compute area for every combination, except:
  i. Also check for point inside polygon - use ray cast algorithm
  ii. Also check that no edges of rectangle intersect (non-parallel) edges of outer polygon
  iii. If either i or ii are not met, discard this area
3. Get largest
*/

func main() {
	mainPart2()
}

func mainPart2() {
	scanner := bufio.NewScanner(os.Stdin)
	points := ReadPoints(scanner)
	lines := MakeLines(points)
	for _, l := range lines {
		fmt.Println(l)
	}
	boxes := MakeBoxes(points)
	areas := make([]int, 0)
	for _, box := range boxes {
		collides := box.InnerBoxCollides(lines)
		if !collides {
			areas = append(areas, box.area)
		}
	}
	maxArea := MaxArea(areas)
	fmt.Printf("Max areas: %v\n", maxArea)
}

func MakeLines(points []Point2D) []LineSegment {
	numPoints := len(points)
	lines := make([]LineSegment, numPoints)
	for i := 0; i < numPoints-1; i++ {
		a := points[i]
		b := points[i+1]
		lines[i] = NewLineSegment(a, b)
	}
	lines[numPoints-1] = NewLineSegment(points[numPoints-1], points[0])
	return lines
}

func NewLineSegment(p1 Point2D, p2 Point2D) LineSegment {
	isHorizontal := p1.y == p2.y

	var a, b Point2D
	if isHorizontal {
		if p1.x < p2.x {
			a = p1
			b = p2
		} else {
			a = p2
			b = p1
		}
	} else {
		if p1.y < p2.y {
			a = p1
			b = p2
		} else {
			a = p2
			b = p1
		}
	}

	return LineSegment{a: a, b: b}
}

func mainPart1() {
	scanner := bufio.NewScanner(os.Stdin)
	points := ReadPoints(scanner)
	areas := ComputeAreas(points)
	for i, a := range areas {
		fmt.Printf("Area %2v: %2v\n", i, a)
	}
	maxArea := MaxArea(areas)
	fmt.Printf("Max area: %v\n", maxArea)
}

type Point2D struct {
	x int
	y int
}

type LineSegment struct {
	a Point2D
	b Point2D
}

func (l LineSegment) IsHorizontal() bool {
	return l.a.y == l.b.y
}

func (l LineSegment) IsVertical() bool {
	return l.a.x == l.b.x
}

func (l LineSegment) Intersects(l2 LineSegment) (bool, bool) {
	/*
		Returns if the lines cross, and if the collision is across a point
		Parallel lines are considered as not colliding, even if overlapping
	*/
	if l.IsHorizontal() && l2.IsHorizontal() {
		sameRow := l.a.y == l2.a.y
		pointOverlap := (l.a.x >= l2.a.x && l.a.x <= l2.b.x) || (l.b.x >= l2.a.x && l.b.x <= l2.b.x)
		return false, sameRow && pointOverlap
	}
	if l.IsVertical() && l2.IsVertical() {
		sameColumn := l.a.x == l2.a.x
		pointOverlap := (l.a.y >= l2.a.y && l.a.y <= l2.b.y) || (l.b.y >= l2.a.y && l.b.y <= l2.b.y)
		return false, sameColumn && pointOverlap
	}

	var horiz, vert LineSegment
	if l.IsHorizontal() {
		// handle l horizontal & l2 vertical case
		horiz = l
		vert = l2
	} else {
		// handle l2 horizontal & l vertical case
		horiz = l2
		vert = l
	}
	// We always examine from the perspective of the horizontal line
	xCrosses := horiz.a.x <= vert.a.x && horiz.b.x >= vert.b.x
	yBetween := false
	if vert.a.y > vert.b.y {
		yBetween = horiz.a.y < vert.a.y && horiz.b.y > vert.b.y
	} else {
		yBetween = horiz.a.y > vert.a.y && horiz.b.y < vert.b.y
	}
	collision := xCrosses && yBetween
	vertPointCollision := horiz.a.y == vert.a.y || horiz.a.y == vert.b.y
	horizPointCollision := horiz.a.x == vert.a.x || horiz.b.x == vert.a.x
	return collision, (vertPointCollision || horizPointCollision)
}

func (l LineSegment) String() string {
	return fmt.Sprintf("LineSegment [a: %v, b: %v]", l.a, l.b)
}

func (p Point2D) Area(o Point2D) int {
	dx := int(math.Abs(float64(p.x) - float64(o.x)))
	dy := int(math.Abs(float64(p.y) - float64(o.y)))
	dx += 1
	dy += 1
	area := dx * dy
	if area < 0 {
		return -area
	}
	return area
}

func (p Point2D) Box(o Point2D) Box {
	area := p.Area(o)

	var upper, lower Point2D
	var topLeft, topRight, bottomLeft, bottomRight Point2D
	if p.y > o.y {
		upper = p
		lower = o
	} else {
		upper = o
		lower = p
	}
	if upper.x < lower.x {
		topLeft = upper
		topRight = Point2D{x: lower.x, y: upper.y}
		bottomLeft = Point2D{x: upper.x, y: lower.y}
		bottomRight = lower
	} else {
		topLeft = Point2D{x: lower.x, y: upper.y}
		topRight = upper
		bottomLeft = lower
		bottomRight = Point2D{x: upper.x, y: lower.y}
	}
	topLine := NewLineSegment(topLeft, topRight)
	rightLine := NewLineSegment(bottomRight, topRight)
	bottomLine := NewLineSegment(bottomLeft, bottomRight)
	leftLine := NewLineSegment(bottomLeft, topLeft)

	return Box{
		a:          p,
		b:          o,
		topSide:    topLine,
		rightSide:  rightLine,
		bottomSide: bottomLine,
		leftSide:   leftLine,
		area:       area,
	}
}

func (b Box) innerBox() Box {
	maxY := b.a.y
	minY := b.b.y
	maxX := b.a.x
	minX := b.b.x
	if b.b.y > maxY {
		maxY = b.b.y
		minY = b.a.y
	}
	if b.b.x > maxX {
		maxX = b.b.x
		minX = b.a.x
	}
	maxY--
	maxX--
	minY++
	minX++
	topLeft := Point2D{x: minX, y: maxY}
	return topLeft.Box(Point2D{x: maxX, y: minY})
}

func (b Box) InnerBoxCollides(lines []LineSegment) bool {
	innerBox := b.innerBox()
	return innerBox.Collides(lines)
}

func (b Box) Collides(lines []LineSegment) bool {
	boxLines := make([]LineSegment, 4)
	// DrawBoxAndLines(b, lines)
	for _, l1 := range lines {
		boxLines[0] = b.topSide
		boxLines[1] = b.leftSide
		boxLines[2] = b.rightSide
		boxLines[3] = b.bottomSide
		for _, l2 := range boxLines {
			collision, pointCollision := l1.Intersects(l2)
			if collision || pointCollision {
				return true
			}
		}
	}
	return false
}

func DrawBoxAndLines(b Box, lines []LineSegment) {
	maxX := 0
	maxY := 0
	fmt.Printf("%v\n", b)
	for _, l := range lines {
		if l.a.x > maxX {
			maxX = l.a.x
		}
		if l.b.x > maxX {
			maxX = l.b.x
		}
		if l.a.y > maxY {
			maxY = l.a.y
		}
		if l.b.y > maxY {
			maxY = l.b.y
		}
	}

	maxX += 2
	maxY += 2
	for j := 0; j < maxY; j++ {
		y := maxY - j
		for i := 0; i < maxX; i++ {
			x := i
			p := Point2D{x: x, y: y}
			char := "."
			for _, l := range lines {
				if p == l.a || p == l.b {
					char = "#"
					break
				}
			}
			boxLines := []LineSegment{b.topSide, b.leftSide, b.bottomSide, b.rightSide}
			for _, l := range boxLines {
				if p == l.a || p == l.b {
					char = "O"
					break
				}
			}
			fmt.Print(char)

		}
		fmt.Println()
	}
}

func (b Box) String() string {
	return fmt.Sprintf("Box [a: %v, b: %v, left: %v, top: %v, right: %v, bottom: %v, area: %v]", b.a, b.b, b.leftSide, b.topSide, b.rightSide, b.bottomSide, b.area)
}

func (p Point2D) String() string {
	return fmt.Sprintf("Point2D [x: %v, y: %v]", p.x, p.y)
}

type Box struct {
	a          Point2D
	b          Point2D
	topSide    LineSegment
	leftSide   LineSegment
	rightSide  LineSegment
	bottomSide LineSegment
	area       int
}

func ReadPoints(scanner *bufio.Scanner) []Point2D {
	points := make([]Point2D, 0)
	for scanner.Scan() {
		line := scanner.Text()
		coordinateStrs := strings.Split(line, ",")
		x, err := strconv.Atoi(coordinateStrs[0])
		if err != nil {
			log.Fatalf("")
		}
		y, err := strconv.Atoi(coordinateStrs[1])
		if err != nil {
			log.Fatalf("")
		}
		points = append(points, Point2D{x: x, y: y})
	}
	return points
}

func MakeBoxes(points []Point2D) []Box {
	numPoints := len(points)
	boxes := make([]Box, 0)
	for i := 0; i < numPoints-1; i++ {
		for j := i + 1; j < numPoints; j++ {
			p1 := points[i]
			p2 := points[j]
			box := p1.Box(p2)
			boxes = append(boxes, box)
		}
	}
	return boxes
}

func ComputeAreas(points []Point2D) []int {
	numPoints := len(points)
	areas := make([]int, 0)
	for i, p1 := range points {
		for j := i + 1; j < numPoints; j++ {
			p2 := points[j]
			areas = append(areas, p1.Area(p2))
		}
	}
	return areas
}

func MaxArea(areas []int) int {
	slices.Sort(areas)
	numAreas := len(areas)
	return areas[numAreas-1]
}
