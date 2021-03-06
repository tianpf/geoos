package space

import (
	"errors"

	"github.com/spatial-go/geoos/algorithm/matrix"
	"github.com/spatial-go/geoos/algorithm/measure"
)

// LineString represents a set of points to be thought of as a polyline.
type LineString matrix.LineMatrix

// GeoJSONType returns the GeoJSON type for the linestring.
func (ls LineString) GeoJSONType() string {
	return TypeLineString
}

// Dimensions returns 1 because a LineString is a 1d object.
func (ls LineString) Dimensions() int {
	return 1
}

// Nums num of linestrings
func (ls LineString) Nums() int {
	return 1
}

// Bound returns a rect around the line string. Uses rectangular coordinates.
func (ls LineString) Bound() Bound {
	if len(ls) == 0 {
		return emptyBound
	}

	b := Bound{ls[0], ls[0]}
	for _, p := range ls {
		b = b.Extend(p)
	}

	return b
}

// EqualLineString compares two line strings. Returns true if lengths are the same
// and all points are Equal.
func (ls LineString) EqualLineString(lineString LineString) bool {
	if len(ls) != len(lineString) {
		return false
	}
	for i, v := range ls.ToPointArray() {
		if !v.Equal(Point(lineString[i])) {
			return false
		}
	}
	return true
}

// Equal checks if the LineString represents the same Geometry or vector.
func (ls LineString) Equal(g Geometry) bool {
	if g.GeoJSONType() != ls.GeoJSONType() {
		return false
	}
	return ls.EqualLineString(g.(LineString))
}

// EqualsExact Returns true if the two Geometrys are exactly equal,
// up to a specified distance tolerance.
// Two Geometries are exactly equal within a distance tolerance
func (ls LineString) EqualsExact(g Geometry, tolerance float64) bool {
	if ls.GeoJSONType() != g.GeoJSONType() {
		return false
	}
	line := g.(LineString)
	if ls.IsEmpty() && g.IsEmpty() {
		return true
	}
	if ls.IsEmpty() != g.IsEmpty() {
		return false
	}
	if len(ls) != len(line) {
		return false
	}

	for i, v := range ls {
		if Point(v).EqualsExact(Point(line[i]), tolerance) {
			return false
		}
	}
	return true
}

// Area returns the area of a polygonal geometry. The area of a LineString is 0.
func (ls LineString) Area() (float64, error) {
	return 0.0, nil
}

// ToPointArray returns the PointArray
func (ls LineString) ToPointArray() (la []Point) {
	for _, v := range ls {
		la = append(la, v)
	}
	return
}

// ToLineArray returns the LineArray
func (ls LineString) ToLineArray() (lines []Line) {
	for i := 0; i < len(ls)-1; i++ {
		lines = append(lines, Line{Point(ls[i]), Point(ls[i+1])})
	}
	return
}

// IsEmpty returns true if the Geometry is empty.
func (ls LineString) IsEmpty() bool {
	return ls == nil || len(ls) == 0
}

// Distance returns distance Between the two Geometry.
func (ls LineString) Distance(g Geometry) (float64, error) {
	elem := &Element{ls}
	return elem.distanceWithFunc(g, measure.PlanarDistance)
}

// SpheroidDistance returns  spheroid distance Between the two Geometry.
func (ls LineString) SpheroidDistance(g Geometry) (float64, error) {
	elem := &Element{ls}
	return elem.distanceWithFunc(g, measure.SpheroidDistance)
}

// Boundary returns the closure of the combinatorial boundary of this space.Geometry.
// The boundary of a lineal geometry is always a zero-dimensional geometry (which may be empty).
func (ls LineString) Boundary() (Geometry, error) {
	if ls.IsClosed() {
		return nil, errors.New("closeedline's boundary should be nil")
	}
	return MultiPoint{ls[0], ls[len(ls)-1]}, nil
}

// IsClosed Returns TRUE if the LINESTRING's start and end points are coincident.
// For Polyhedral Surfaces, reports if the surface is areal (open) or IsC (closed).
func (ls LineString) IsClosed() bool {
	if Point(ls[0]).Equal(Point(ls[len(ls)-1])) {
		return true
	}
	return false
}

// Length Returns the length of this LineString
func (ls LineString) Length() float64 {
	return measure.OfLine(matrix.LineMatrix(ls))
}

// IsSimple returns true if this space.Geometry has no anomalous geometric points,
// such as self intersection or self tangency.
func (ls LineString) IsSimple() bool {
	elem := ElementValid{ls}
	return elem.IsSimple()
}
