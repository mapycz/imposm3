package test

import (
	"database/sql"
	"io/ioutil"

	"testing"

	"github.com/omniscale/imposm3/geom/geos"
)

func TestMultiLineString(t *testing.T) {
	if testing.Short() {
		t.Skip("system test skipped with -test.short")
	}
	t.Parallel()

	ts := importTestSuite{
		name: "multilinestring",
	}

	t.Run("Prepare", func(t *testing.T) {
		var err error

		ts.dir, err = ioutil.TempDir("", "imposm_test")
		if err != nil {
			t.Fatal(err)
		}
		ts.config = importConfig{
			connection:      "postgis://",
			cacheDir:        ts.dir,
			osmFileName:     "build/multilinestring.pbf",
			mappingFileName: "multilinestring_mapping.yml",
		}
		ts.g = geos.NewGeos()

		ts.db, err = sql.Open("postgres", "sslmode=disable")
		if err != nil {
			t.Fatal(err)
		}
		ts.dropSchemas()
	})

	const mlsTable = "osm_multilinestring"

	t.Run("Import", func(t *testing.T) {
		if ts.tableExists(t, ts.dbschemaImport(), mlsTable) != false {
			t.Fatalf("table %s exists in schema %s", mlsTable, ts.dbschemaImport())
		}
		ts.importOsm(t)
		if ts.tableExists(t, ts.dbschemaImport(), mlsTable) != true {
			t.Fatalf("table %s does not exists in schema %s", mlsTable, ts.dbschemaImport())
		}
	})

	t.Run("Deploy", func(t *testing.T) {
		ts.deployOsm(t)
		if ts.tableExists(t, ts.dbschemaImport(), mlsTable) != false {
			t.Fatalf("table %s exists in schema %s", mlsTable, ts.dbschemaImport())
		}
		if ts.tableExists(t, ts.dbschemaProduction(), mlsTable) != true {
			t.Fatalf("table %s does not exists in schema %s", mlsTable, ts.dbschemaProduction())
		}
	})

	t.Run("CheckMultiLineStringGeometry", func(t *testing.T) {
		element := checkElem{mlsTable, -1077494, "*", nil}
		ts.assertGeomType(t, element, "MultiLineString")
		ts.assertGeomValid(t, element)
		ts.assertGeomLength(t, element, 22005)
	})

	t.Run("CheckLineStringGeometry", func(t *testing.T) {
		element := checkElem{mlsTable, 40881141, "*", nil}
		ts.assertGeomType(t, element, "LineString")
		ts.assertGeomValid(t, element)
		ts.assertGeomLength(t, element, 207)
	})

	t.Run("CheckFilters", func(t *testing.T) {
		if records := ts.queryRows(t, mlsTable, 25948796); len(records) > 0 {
			t.Fatalf("The way 25948796 should be filtered out by typed filter")
		}
		if records := ts.queryRows(t, mlsTable, 27037886); len(records) > 0 {
			t.Fatalf("The way 27037886 should be filtered out as it is closed path with area=yes")
		}
	})

	t.Run("RelationTypesFilter", func(t *testing.T) {
		if records := ts.queryRows(t, "osm_multilinestring_no_relations", -1077494); len(records) > 0 {
			t.Fatalf("The relation -1077494 should not be imported due to empty relation_types")
		}
	})

	t.Run("Update", func(t *testing.T) {
		ts.updateOsm(t, "build/multilinestring.osc.gz")
	})

	t.Run("CheckFilters2", func(t *testing.T) {
		if records := ts.queryRows(t, mlsTable, 27037886); len(records) == 0 {
			t.Fatalf("The way 27037886 should now be there as we removed area=yes in the update")
		}
	})

	t.Run("CheckNewRelation", func(t *testing.T) {
		if records := ts.queryRows(t, mlsTable, -8747972); len(records) == 0 {
			t.Fatalf("The relation -8747972 should be created")
		}
	})
}

