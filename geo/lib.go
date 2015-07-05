package geo

import (
	"fmt"
	"math"

	"github.com/kevinburke/rct/tracks"
)

// A Point is a representation of (x, y, z) coordinates
type Point [3]float64

func Render(elems []tracks.Element) ([]Point, []Point) {
	current := Point{0, 0, 0}
	left := []Point{current}
	var direction tracks.DirectionDelta
	direction = 0
	for _, elem := range elems {
		segm := tracks.TS_MAP[elem.Segment.Type]
		fmt.Println(segm.Type)
		p := Diff(segm, direction)
		fmt.Println(p[2])
		current[0] += p[0]
		current[1] += p[1]
		current[2] += p[2]
		fmt.Println(current[2])
		left = append(left, current)
		direction += segm.DirectionDelta
		for direction > 360 {
			direction -= 360
		}
	}
	return left, []Point{}
}

// Round to the nearest integer
func round(f float64) int {
	return int(math.Floor(f + 0.5))
}

// Advance all of the values by one track segment.
func advanceTrack(ts *tracks.Segment, ΔE int, ΔForward int, ΔSideways int,
	direction tracks.DirectionDelta) (int, int, int, tracks.DirectionDelta) {

	// XXX
	ΔE += 0

	fdirection := float64(direction)
	ΔForward += round(cosdeg(fdirection) * float64(ts.ForwardDelta))
	ΔForward += round(sindeg(fdirection) * float64(ts.SidewaysDelta))

	ΔSideways += round(sindeg(fdirection) * float64(ts.ForwardDelta))
	ΔSideways += round(cosdeg(fdirection) * float64(ts.SidewaysDelta))

	direction += ts.DirectionDelta
	for ; direction >= 360; direction -= 360 {
	}

	return ΔE, ΔForward, ΔSideways, direction
}

// sindeg computes sines in degrees
func sindeg(deg float64) float64 {
	for ; deg >= 360; deg -= 360 {
	}
	if round(deg)%180 == 0 {
		return 0
	} else if round(deg) == 90 {
		return 1
	} else if round(deg) == 270 {
		return -1
	} else {
		return math.Sin(deg * math.Pi / 180)
	}
}

// computes sines in degrees
func cosdeg(deg float64) float64 {
	for ; deg >= 360; deg -= 360 {
	}
	if round(deg) == 0 {
		return 1
	} else if round(deg) == 90 || round(deg) == 270 {
		return 0
	} else if round(deg) == 180 {
		return -1
	} else {
		return math.Sin(deg * math.Pi / 180)
	}
}

// PositionChange takes a travel direction in degrees and a track segment and
// returns the distance traveled in the X, Y, and Z directions.
func Diff(ts *tracks.Segment, direction tracks.DirectionDelta) *Point {
	//rotate around the Z axis: http://stackoverflow.com/a/14609567/329700
	x := float64(ts.ForwardDelta)*cosdeg(float64(direction)) - float64(ts.SidewaysDelta)*sindeg(float64(direction))
	y := float64(ts.ForwardDelta)*sindeg(float64(direction)) + float64(ts.SidewaysDelta)*cosdeg(float64(direction))
	return &Point{x, y, float64(ts.ElevationDelta)}
}

// Test whether the ride's track forms a continuous circuit. Does not test
// whether the ride collides with itself.
func IsCircuit(t *tracks.Data) bool {
	// X and Y don't really make sense as variable names, easier to think about
	// relative changes
	eΔ, forwardΔ, sidewaysΔ := 0, 0, 0
	direction := tracks.DIR_STRAIGHT
	if len(t.Elements) == 0 {
		return false
	}
	for i := range t.Elements {
		ts := t.Elements[i].Segment
		eΔ, forwardΔ, sidewaysΔ, direction = advanceTrack(
			ts, eΔ, forwardΔ, sidewaysΔ, direction)
	}
	return forwardΔ == 0 && sidewaysΔ == 0 && eΔ == 0
}

// Detect whether the track collides with itself.
func HasCollision(t *tracks.Data) bool {
	matrix := make([][][]bool, 100)
	for i := range matrix {
		matrix[i] = make([][]bool, 100)
		for j := range matrix[i] {
			matrix[i][j] = make([]bool, 100)
		}
	}
	eΔ, forwardΔ, sidewaysΔ := 0, 0, 0
	direction := tracks.DIR_STRAIGHT
	for i := range t.Elements {
		ts := t.Elements[i].Segment
		eΔ, forwardΔ, sidewaysΔ, direction = advanceTrack(
			ts, eΔ, forwardΔ, sidewaysΔ, direction)
		// if there already exists a piece there, we can't build.
		if matrix[forwardΔ][sidewaysΔ][eΔ] {
			return true
		}
		matrix[forwardΔ][sidewaysΔ][eΔ] = true
	}
	return false
}