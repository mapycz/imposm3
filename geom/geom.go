package geom

import (
	"goposm/element"
	"goposm/geom/geos"
)

func PointWKB(g *geos.GEOS, node element.Node) ([]byte, error) {
	coordSeq, err := g.CreateCoordSeq(1, 2)
	if err != nil {
		return nil, err
	}
	// coordSeq inherited by LineString
	coordSeq.SetXY(g, 0, node.Long, node.Lat)
	geom, err := coordSeq.AsPoint(g)
	if err != nil {
		return nil, err
	}
	defer g.Destroy(geom)
	return g.AsWKB(geom)
}

func LineStringWKB(g *geos.GEOS, nodes []element.Node) ([]byte, error) {
	coordSeq, err := g.CreateCoordSeq(uint32(len(nodes)), 2)
	if err != nil {
		return nil, err
	}
	// coordSeq inherited by LineString
	for i, nd := range nodes {
		coordSeq.SetXY(g, uint32(i), nd.Long, nd.Lat)
	}
	geom, err := coordSeq.AsLineString(g)
	if err != nil {
		return nil, err
	}
	defer g.Destroy(geom)
	return g.AsWKB(geom)
}

func PolygonWKB(g *geos.GEOS, nodes []element.Node) ([]byte, error) {
	coordSeq, err := g.CreateCoordSeq(uint32(len(nodes)), 2)
	if err != nil {
		return nil, err
	}
	// coordSeq inherited by LineString, no destroy
	for i, nd := range nodes {
		err := coordSeq.SetXY(g, uint32(i), nd.Long, nd.Lat)
		if err != nil {
			return nil, err
		}
	}
	geom, err := coordSeq.AsLinearRing(g)
	if err != nil {
		return nil, err
	}
	// geom inherited by Polygon, no destroy
	geom = g.CreatePolygon(geom, nil)
	if err != nil {
		return nil, err
	}
	defer g.Destroy(geom)

	return g.AsWKB(geom)
}
