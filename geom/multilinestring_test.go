package geom

import (
	"testing"

	"github.com/omniscale/imposm3/element"
	"github.com/omniscale/imposm3/geom/geos"
)

func TestSimpleMultiLineString(t *testing.T) {
	w1 := makeWay(1, element.Tags{}, []coord{
		{1, 1, 0},
		{2, 2, 0},
	})
	w2 := makeWay(2, element.Tags{}, []coord{
		{3, 2, 0},
		{4, 3, 0},
	})

	rel := element.Relation{
		OSMElem: element.OSMElem{Id: 1, Tags: element.Tags{}}}
	rel.Members = []element.Member{
		{Id: 1, Type: element.WAY, Role: "", Way: &w1},
		{Id: 2, Type: element.WAY, Role: "", Way: &w2},
	}

	geom, err := BuildMultiLinestring(&rel, 3857)

	if err != nil {
		t.Fatal(err)
	}

	g := geos.NewGeos()
	defer g.Finish()

	if !g.IsValid(geom) {
		t.Fatal("geometry not valid", g.AsWkt(geom))
	}

	if length := geom.Length(); length != 2 {
		t.Fatal("length invalid", length)
	}
}
