package planar

import (
	"reflect"
	"testing"

	"github.com/spatial-go/geoos/encoding/wkt"
	"github.com/spatial-go/geoos/space"
)

func TestAlgorithm_Area(t *testing.T) {
	const polygon = `POLYGON((-1 -1, 1 -1, 1 1, -1 1, -1 -1))`
	geometry, _ := wkt.UnmarshalString(polygon)
	type args struct {
		g space.Geometry
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{name: "area", args: args{g: geometry}, want: 4.0, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalStrategy().Area(tt.args.g)
			if (err != nil) != tt.wantErr {
				t.Errorf("Area() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Area() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlgorithm_Length(t *testing.T) {
	const line = `LINESTRING(0 0, 1 1)`
	geometry, _ := wkt.UnmarshalString(line)
	type args struct {
		g space.Geometry
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{name: "lengh", args: args{g: geometry}, want: 1.4142135623730951, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			G := NormalStrategy()
			got, err := G.Length(tt.args.g)
			if (err != nil) != tt.wantErr {
				t.Errorf("Length() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Length() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlgorithm_Boundary(t *testing.T) {
	const sourceLine = `LINESTRING(1 1,0 0, -1 1)`
	const expectLine = `MULTIPOINT(1 1,-1 1)`
	sLine, _ := wkt.UnmarshalString(sourceLine)
	eLine, _ := wkt.UnmarshalString(expectLine)

	const sourcePolygon = `POLYGON((1 1,0 0, -1 1, 1 1))`
	const expectPolygon = `LINESTRING(1 1,0 0,-1 1,1 1)`
	sPolygon, _ := wkt.UnmarshalString(sourcePolygon)
	ePolygon, _ := wkt.UnmarshalString(expectPolygon)

	// const multiPolygon = `POLYGON (( 10 130, 50 190, 110 190, 140 150, 150 80, 100 10, 20 40, 10 130 ),
	// ( 70 40, 100 50, 120 80, 80 110, 50 90, 70 40 ))`
	// const expectMultiPolygon = `MULTILINESTRING((10 130,50 190,110 190,140 150,150 80,100 10,20 40,10 130),
	// (70 40,100 50,120 80,80 110,50 90,70 40))`

	smultiPolygon, _ := wkt.UnmarshalString(sourceLine)
	emultiPolygon, _ := wkt.UnmarshalString(expectLine)

	type args struct {
		g space.Geometry
	}
	tests := []struct {
		name    string
		args    args
		want    space.Geometry
		wantErr bool
	}{
		{name: "line", args: args{g: sLine}, want: eLine, wantErr: false},
		{name: "polygon", args: args{g: sPolygon}, want: ePolygon, wantErr: false},
		{name: "multiPolygon", args: args{g: smultiPolygon}, want: emultiPolygon, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			G := NormalStrategy()
			got, err := G.Boundary(tt.args.g)
			if (err != nil) != tt.wantErr {
				t.Errorf("Boundary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(wkt.MarshalString(got))
			t.Log(wkt.MarshalString(tt.want))
			if !got.Equal(tt.want) {
				t.Errorf("Boundary() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlgorithm_UnaryUnion(t *testing.T) {
	multiPolygon, _ := wkt.UnmarshalString(`MULTIPOLYGON(((0 0, 10 0, 10 10, 0 10, 0 0)), ((5 5, 15 5, 15 15, 5 15, 5 5)))`)
	expectPolygon, _ := wkt.UnmarshalString(`POLYGON((5 10,0 10,0 0,10 0,10 5,15 5,15 15,5 15,5 10))`)
	type args struct {
		g space.Geometry
	}
	tests := []struct {
		name    string
		G       GEOAlgorithm
		args    args
		want    []space.Geometry
		wantErr bool
	}{
		{name: "UnaryUnion Polygon", args: args{g: multiPolygon}, want: []space.Geometry{expectPolygon}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			G := NormalStrategy()
			gotGeometry, err := G.UnaryUnion(tt.args.g)

			if (err != nil) != tt.wantErr {
				t.Errorf("Algorithm UnaryUnion error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			isEqual, _ := G.EqualsExact(gotGeometry, tt.want[0], 0.000001)
			if !isEqual {
				t.Errorf("Algorithm UnaryUnion = %v, want %v", wkt.MarshalString(gotGeometry), wkt.MarshalString(tt.want[0]))
			}
		})
	}
}

func TestAlgorithm_Union(t *testing.T) {
	point01, _ := wkt.UnmarshalString(`POINT(1 2)`)
	point02, _ := wkt.UnmarshalString(`POINT(-2 3)`)
	expectMultiPoint, _ := wkt.UnmarshalString(`MULTIPOINT(1 2,-2 3)`)

	type args struct {
		g1 space.Geometry
		g2 space.Geometry
	}
	tests := []struct {
		name    string
		args    args
		want    space.Geometry
		wantErr bool
	}{
		{name: "union", args: args{g1: point01, g2: point02}, want: expectMultiPoint},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			G := NormalStrategy()
			got, err := G.Union(tt.args.g1, tt.args.g2)
			if (err != nil) != tt.wantErr {
				t.Errorf("Union() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Union() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlgorithm_Distance(t *testing.T) {
	point01, _ := wkt.UnmarshalString(`POINT(1 3)`)
	point02, _ := wkt.UnmarshalString(`POINT(4 7)`)

	type args struct {
		g1 space.Geometry
		g2 space.Geometry
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{name: "distance", args: args{g1: point01, g2: point02}, want: 5, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			G := NormalStrategy()
			got, err := G.Distance(tt.args.g1, tt.args.g2)
			if (err != nil) != tt.wantErr {
				t.Errorf("Distance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Distance() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlgorithm_SphericalDistance(t *testing.T) {
	point01 := space.Point{116.397439, 39.909177}
	point02 := space.Point{116.397725, 39.903079}
	point03 := space.Point{118.1487, 39.586671}

	type args struct {
		p1 space.Point
		p2 space.Point
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{name: "SphericalDistance", args: args{p1: point01, p2: point02}, want: 678.5053586786567, wantErr: false},
		{name: "SphericalDistance", args: args{p1: point01, p2: point03}, want: 153953.98145619757, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			G := NormalStrategy()
			got, _ := G.SphericalDistance(tt.args.p1, tt.args.p2)
			if got != tt.want {
				t.Errorf("SphericalDistance() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlgorithm_IsClosed(t *testing.T) {
	const linestring = `LINESTRING(1 1,2 3,3 2,1 2)`
	const closedLinestring = `LINESTRING(1 1,2 3,3 2,1 2,1 1)`
	line, _ := wkt.UnmarshalString(linestring)
	closedLine, _ := wkt.UnmarshalString(closedLinestring)

	type args struct {
		g space.Geometry
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "line", args: args{g: line}, want: false, wantErr: false},
		{name: "closedLine", args: args{g: closedLine}, want: true, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			G := NormalStrategy()
			got, err := G.IsClosed(tt.args.g)
			t.Log(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("geometry: %s IsClose() error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("geometry: %s IsClose() got = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestAlgorithm_EqualsExact(t *testing.T) {
	geometry1, _ := wkt.UnmarshalString("POINT(116.309878625564 40.0427783817455)")
	geometry2, _ := wkt.UnmarshalString("POINT(116.309878725564 40.0427783827455)")
	geometry3, _ := wkt.UnmarshalString("POINT(116.309877625564 40.0427783827455)")
	type args struct {
		g1        space.Geometry
		g2        space.Geometry
		tolerance float64
	}
	tests := []struct {
		name    string
		G       GEOAlgorithm
		args    args
		want    bool
		wantErr bool
	}{
		{name: "equals exact", args: args{g1: geometry1, g2: geometry2, tolerance: 0.000001}, want: true, wantErr: false},
		{name: "not equals exact", args: args{g1: geometry1, g2: geometry3, tolerance: 0.000001}, want: false, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			G := NormalStrategy()
			got, err := G.EqualsExact(tt.args.g1, tt.args.g2, tt.args.tolerance)
			if (err != nil) != tt.wantErr {
				t.Errorf("GEOAlgorithm.EqualsExact() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GEOAlgorithm.EqualsExact() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlgorithm_Equals(t *testing.T) {
	geometry1, _ := wkt.UnmarshalString("POINT(116.309878625564 40.0427783817455)")
	geometry2, _ := wkt.UnmarshalString("POINT(116.309878725564 40.0427783827455)")
	geometry3, _ := wkt.UnmarshalString("POINT(116.309878625564 40.0427783817455)")
	type args struct {
		g1        space.Geometry
		g2        space.Geometry
		tolerance float64
	}
	tests := []struct {
		name    string
		G       GEOAlgorithm
		args    args
		want    bool
		wantErr bool
	}{
		{name: "equals exact", args: args{g1: geometry1, g2: geometry3}, want: true, wantErr: false},
		{name: "not equals exact", args: args{g1: geometry1, g2: geometry2}, want: false, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			G := NormalStrategy()
			got, err := G.Equals(tt.args.g1, tt.args.g2)
			if (err != nil) != tt.wantErr {
				t.Errorf("GEOAlgorithm.Equals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GEOAlgorithm.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlgorithm_Envelope(t *testing.T) {
	point, _ := wkt.UnmarshalString(`POINT(1 3)`)
	expectPoint, _ := wkt.UnmarshalString(`POINT(1 3)`)

	line, _ := wkt.UnmarshalString(`LINESTRING(0 0, 1 3)`)
	expectPolygon0, _ := wkt.UnmarshalString(`POLYGON((0 0,1 0,1 3,0 3,0 0))`)

	polygon1, _ := wkt.UnmarshalString(`POLYGON((0 0, 0 1, 1.0000001 1, 1.0000001 0, 0 0))`)
	expectPolygon1, _ := wkt.UnmarshalString(`POLYGON((0 0,1.0000001 0,1.0000001 1,0 1,0 0))`)

	polygon2, _ := wkt.UnmarshalString(`POLYGON((0 0, 0 1, 1.0000000001 1, 1.0000000001 0, 0 0))`)
	expectPolygon2, _ := wkt.UnmarshalString(`POLYGON((0 0,1.0000000001 0,1.0000000001 1,0 1,0 0))`)

	type args struct {
		g space.Geometry
	}

	tests := []struct {
		name    string
		args    args
		want    space.Geometry
		wantErr bool
	}{
		{name: "envelope Point", args: args{g: point}, want: expectPoint, wantErr: false},
		{name: "envelope LineString", args: args{g: line}, want: expectPolygon0, wantErr: false},
		{name: "envelope Polygon", args: args{g: polygon1}, want: expectPolygon1, wantErr: false},
		{name: "envelope Polygon", args: args{g: polygon2}, want: expectPolygon2, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			G := NormalStrategy()
			gotGeometry, err := G.Envelope(tt.args.g)

			if (err != nil) != tt.wantErr {
				t.Errorf("GEOAlgorithm.EqualsExact() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			isEqual, _ := G.EqualsExact(gotGeometry, tt.want, 0.000001)
			if !isEqual {
				t.Errorf("GEOAlgorithm.Envelope() = %v, want %v", wkt.MarshalString(gotGeometry), wkt.MarshalString(tt.want))
			}
		})
	}
}

func TestAlgorithm_IsSimple(t *testing.T) {
	const polygon = `POLYGON((1 2, 3 4, 5 6, 5 3, 1 2))`
	const linestring = `LINESTRING(1 1,2 2,2 3.5,1 3,1 2,2 1)`
	poly, _ := wkt.UnmarshalString(polygon)
	line, _ := wkt.UnmarshalString(linestring)

	type args struct {
		g space.Geometry
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "polygon", args: args{g: poly}, want: true, wantErr: false},
		{name: "line", args: args{g: line}, want: false, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			G := NormalStrategy()
			got, err := G.IsSimple(tt.args.g)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsSimple() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsSimple() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlgorithm_IsRing(t *testing.T) {
	const linestring1 = `LINESTRING((1 2, 3 4, 5 6, 5 3, 1 2))`
	const linestring2 = `LINESTRING(1 1,2 2,2 3.5,1 3,1 2,2 1)`
	line1, _ := wkt.UnmarshalString(linestring1)
	line2, _ := wkt.UnmarshalString(linestring2)

	type args struct {
		g space.Geometry
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "polygon", args: args{g: line1}, want: true, wantErr: false},
		{name: "line", args: args{g: line2}, want: false, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			G := NormalStrategy()
			got, err := G.IsSimple(tt.args.g)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsRing() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsRing() got = %v, want %v", got, tt.want)
			}
		})
	}
}
