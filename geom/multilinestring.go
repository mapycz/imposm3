package geom

import (
	"errors"
	"github.com/omniscale/imposm3/element"
	"github.com/omniscale/imposm3/geom/geos"
	"runtime"
)

func BuildMultiLinestring(rel *element.Relation, srid int) (*geos.Geom, error) {
	g := geos.NewGeos()
	g.SetHandleSrid(srid)
	defer g.Finish()

	var lines []*geos.Geom

	for _, member := range rel.Members {
		if member.Way == nil {
			continue
		}

		line, err := LineString(g, member.Way.Nodes)

		// Clear the finalizer created in LineString()
		// as we want to make the object a part of MultiLineString.
		runtime.SetFinalizer(line, nil)

		if err != nil {
			return nil, err
		}

		lines = append(lines, line)
	}

	result := g.MultiLineString(lines)
	if result == nil {
		return nil, errors.New("Error while building multi-linestring.")
	}

	g.DestroyLater(result)

	return result, nil
}
