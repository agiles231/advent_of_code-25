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
1. Get input -> make points
2. Compute distances
3. Select shortest distance: make edge between the two points
    i. Check if either point is in a graph
   	a. 1 point in graph -> Add other to graph
	b. 2 points in graph, merge graphs
	    i. Update existing mappings of all points to the removed graph
	c. 0 points in graph, make a new graph
    ii. Once added to graph, map both points to new graph
4. Using edges made, determine N largest graphs
*/

func main() {
	limitStr := ""
	limit := -1
	var err error
	if len(os.Args) > 1 {
		limitStr = os.Args[1]
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			log.Fatalf("Failure: %v\n", err)
		}
	}
	fmt.Println("Hello, world!")
	scanner := bufio.NewScanner(os.Stdin)
	points, err := ReadPoints(scanner)
	if err != nil {
		log.Fatalf("Failure: %v\n", err)
	}
	PrintPoints(points)
	distances := ComputeDistances(points)

	slices.SortFunc(distances, func(a, b Point3DDistance) int {
		dist := a.distance - b.distance
		if dist > 0 {
			return 1
		} else if dist < 0 {
			return -1
		} else {
			return 0
		}
	})
	PrintDistances(distances)
	graphMap := make(map[Point3D]*Point3DGraph)
	graphs := make([]*Point3DGraph, len(points))

	for i, p := range points {
		g := NewPoint3DGraph(nil, p, 1)
		graphMap[p] = g
		graphs[i] = g
	}

	shortestDistances := distances[:limit]
	for _, d := range shortestDistances {
		g1, ok := graphMap[d.a]
		if !ok {
			log.Fatalf("Failure: Point %v wasn't found in graphs\n", d.a)
		}
		g2, ok := graphMap[d.b]
		if !ok {
			log.Fatalf("Failure: Point %v wasn't found in graphs\n", d.a)
		}
		g1.Merge(g2)
	}

	slices.SortFunc(graphs, func(a, b *Point3DGraph) int {
		return b.size - a.size
	})
	biggestGraphs := graphs[:3]

	computedValue := biggestGraphs[0].size * biggestGraphs[1].size * biggestGraphs[2].size

	fmt.Printf("Value: %v\n", computedValue)

	for _, d := range distances {
		g1, ok := graphMap[d.a]
		if !ok {
			log.Fatalf("Failure: Point %v wasn't found in graphs\n", d.a)
		}
		g2, ok := graphMap[d.b]
		if !ok {
			log.Fatalf("Failure: Point %v wasn't found in graphs\n", d.a)
		}
		g1.Merge(g2)
		if oneSingleGraph(graphs) {
			fmt.Printf("Last 2 X values: %v, %v\n", d.a.x, d.b.x)
			fmt.Printf("Product: %v\n", d.a.x*d.b.x)
			break
		} else {
			fmt.Printf("Still going...\n")
		}

	}
	if !oneSingleGraph(graphs) {
		fmt.Printf("What happened??\n")
	}
}

func oneSingleGraph(graphs []*Point3DGraph) bool {
	root := graphs[0].GetRoot()
	for _, g := range graphs {
		if g.GetRoot() != root {
			return false
		}
	}
	return true
}

func ComputeDistances(points []Point3D) []Point3DDistance {
	distances := make([]Point3DDistance, 0)
	for i, p1 := range points {
		for j := i + 1; j < len(points); j++ {
			p2 := points[j]
			distances = append(distances, p1.Distance(p2))
		}
	}
	return distances
}

func PrintDistances(dists []Point3DDistance) {
	for i, d := range dists {
		fmt.Printf("Distance %2v: %5v\n", i, d.String())
	}
}

func PrintPoints(points []Point3D) {
	for i, p := range points {
		fmt.Printf("Point %2v: %v\n", i, p.String())
	}
}

func ReadPoints(scanner *bufio.Scanner) ([]Point3D, error) {
	points := make([]Point3D, 0)
	for scanner.Scan() {
		line := scanner.Text()
		numStrs := strings.Split(line, ",")
		nums, err := ConvertToNums(numStrs)
		if err != nil {
			return nil, err
		}
		points = append(points, NewPoint3D(nums[0], nums[1], nums[2]))
	}
	return points, nil
}

func ConvertToNums(numStrs []string) ([]int, error) {
	nums := make([]int, 0)
	for _, numStr := range numStrs {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return nil, err
		}
		nums = append(nums, num)
	}
	return nums, nil
}

type Point3D struct {
	x int
	y int
	z int
}

func NewPoint3D(x, y, z int) Point3D {
	return Point3D{
		x: x,
		y: y,
		z: z,
	}
}

func (p Point3D) String() string {
	return fmt.Sprintf("Point3D [x: %8v, y: %8v, z: %8v]\n", p.x, p.y, p.z)
}

func (p Point3D) Distance(o Point3D) Point3DDistance {
	dx := p.x - o.x
	dy := p.y - o.y
	dz := p.z - o.z
	dxSquared := dx * dx
	dySquared := dy * dy
	dzSquared := dz * dz
	return Point3DDistance{
		a:        p,
		b:        o,
		distance: math.Sqrt(float64(dxSquared) + float64(dySquared) + float64(dzSquared)),
	}
}

type Point3DDistance struct {
	a        Point3D
	b        Point3D
	distance float64
}

func (p Point3DDistance) String() string {
	return fmt.Sprintf("Point3DDistance [A: %v, B: %v, Distance: %10v \n]", p.a.String(), p.b.String(), p.distance)
}

type Point3DGraph struct {
	root  *Point3DGraph
	point Point3D
	size  int
}

func NewPoint3DGraph(root *Point3DGraph, point Point3D, size int) *Point3DGraph {
	g := &Point3DGraph{
		root:  root,
		point: point,
		size:  size,
	}
	if g.root == nil {
		g.root = g
	}
	return g
}

func (g *Point3DGraph) GetRoot() *Point3DGraph {
	if g.root == g {
		return g
	} else {
		return g.root.GetRoot()
	}
}

func (g *Point3DGraph) Merge(o *Point3DGraph) {
	rg := g.GetRoot()
	ro := o.GetRoot()
	if ro == rg {
		return
	}
	if rg.size < ro.size {
		temp := ro
		ro = rg
		rg = temp
	}
	rg.size += ro.root.size
	ro.root = rg
}
