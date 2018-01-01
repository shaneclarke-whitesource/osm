package osm

import (
	"bytes"
	"context"
	"encoding/xml"
	"reflect"
	"testing"
)

func TestChange(t *testing.T) {
	data := []byte(`
<osmChange version="0.6" generator="OpenStreetMap server" copyright="OpenStreetMap and contributors" attribution="http://www.openstreetmap.org/copyright" license="http://opendatacommons.org/licenses/odbl/1-0/">
<create>
<node id="2780675158" changeset="21598503" timestamp="2014-04-10T00:43:05Z" version="1" visible="true" user="jeroenrozema74" uid="183483" lat="50.7107023" lon="6.0043943"/>
</create>
<create>
<node id="2780675159" changeset="21598503" timestamp="2014-04-10T00:43:05Z" version="1" visible="true" user="jeroenrozema74" uid="183483" lat="50.710755" lon="5.9998612"/>
</create>
<create>
<way id="273193870" changeset="21598503" timestamp="2014-04-10T00:43:07Z" version="1" visible="true" user="jeroenrozema74" uid="183483">
<nd ref="2780675158"/>
<nd ref="2780675160"/>
<nd ref="2780675161"/>
<nd ref="2780675162"/>
<nd ref="2780675164"/>
<tag k="highway" v="unclassified"/>
</way>
</create>
<modify>
<way id="24830559" changeset="21598503" timestamp="2014-04-10T00:43:07Z" version="9" visible="true" user="jeroenrozema74" uid="183483">
<nd ref="269419649"/>
<nd ref="269419627"/>
<nd ref="166810716"/>
<nd ref="1149226084"/>
<nd ref="269704932"/>
<nd ref="269419651"/>
<nd ref="2751084538"/>
<nd ref="269419653"/>
<nd ref="269419654"/>
<nd ref="269419655"/>
<nd ref="2780675158"/>
<nd ref="269658287"/>
<nd ref="2351330343"/>
<nd ref="269419658"/>
<tag k="highway" v="tertiary"/>
<tag k="name" v="Krikelsteinstraße"/>
</way>
</modify>
<delete>
<way id="252107750" changeset="21598503" timestamp="2014-04-10T00:43:11Z" version="3" visible="false" user="jeroenrozema74" uid="183483"/>
</delete>
<delete>
<way id="252107748" changeset="21598503" timestamp="2014-04-10T00:43:11Z" version="4" visible="false" user="jeroenrozema74" uid="183483"/>
</delete>
<delete>
<node id="301847601" changeset="21598503" timestamp="2014-04-10T00:43:15Z" version="2" visible="false" user="jeroenrozema74" uid="183483"/>
</delete>
</osmChange>
	`)

	c := Change{}
	err := xml.Unmarshal(data, &c)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if l := len(c.Create.Nodes); l != 2 {
		t.Errorf("incorrect number of created nodes, got %v", l)
	}

	if l := len(c.Create.Ways); l != 1 {
		t.Errorf("incorrect number of created ways, got %v", l)
	}

	if v := c.Create.Nodes[0].ID; v != 2780675158 {
		t.Errorf("incorrect node id, got %v", v)
	}

	if v := c.Create.Nodes[1].ID; v != 2780675159 {
		t.Errorf("incorrect node id, got %v", v)
	}

	if v := c.Create.Ways[0].ID; v != 273193870 {
		t.Errorf("incorrect way id, got %v", v)
	}

	if l := len(c.Modify.Ways); l != 1 {
		t.Errorf("incorrect number of modified ways, got %v", l)
	}

	if v := c.Modify.Ways[0].ID; v != 24830559 {
		t.Errorf("incorrect way id, got %v", v)
	}

	if l := len(c.Delete.Nodes); l != 1 {
		t.Errorf("incorrect number of deleted nodes, got %v", l)
	}

	if l := len(c.Delete.Ways); l != 2 {
		t.Errorf("incorrect number of deleted ways, got %v", l)
	}

	if v := c.Delete.Ways[0].ID; v != 252107750 {
		t.Errorf("incorrect way id, got %v", v)
	}

	if v := c.Delete.Ways[1].ID; v != 252107748 {
		t.Errorf("incorrect way id, got %v", v)
	}

	if v := c.Delete.Nodes[0].ID; v != 301847601 {
		t.Errorf("incorrect node id, got %v", v)
	}

	// empty change
	data = []byte(`<osmChange version="0.6" generator="OpenStreetMap server" copyright="OpenStreetMap and contributors" attribution="http://www.openstreetmap.org/copyright" license="http://opendatacommons.org/licenses/odbl/1-0/"/>`)

	c = Change{}
	err = xml.Unmarshal(data, &c)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if c.Create != nil {
		t.Errorf("create should be nil for empty change")
	}

	if c.Modify != nil {
		t.Errorf("modify should be nil for empty change")
	}

	if c.Delete != nil {
		t.Errorf("delete should be nil for empty change")
	}
}

func TestChange_MarshalXML(t *testing.T) {
	// correct case of name
	c := Change{
		Version:     0.6,
		Generator:   "osm-go",
		Copyright:   "copyright1",
		Attribution: "attribution1",
		License:     "license1",
		Create: &OSM{
			Nodes: Nodes{
				&Node{ID: 123},
			},
		},
	}

	data, err := xml.Marshal(c)
	if err != nil {
		t.Fatalf("xml marshal error: %v", err)
	}

	expected := `<osmChange version="0.6" generator="osm-go" copyright="copyright1" attribution="attribution1" license="license1"><create><node id="123" lat="0" lon="0" user="" uid="0" visible="false" version="0" changeset="0" timestamp="0001-01-01T00:00:00Z"></node></create></osmChange>`
	if !bytes.Equal(data, []byte(expected)) {
		t.Errorf("incorrect marshal, got: %s", string(data))
	}

	// omit attributes if not defined
	c = Change{}
	data, err = xml.Marshal(c)
	if err != nil {
		t.Fatalf("xml marshal error: %v", err)
	}

	expected = `<osmChange></osmChange>`
	if !bytes.Equal(data, []byte(expected)) {
		t.Errorf("incorrect marshal, got: %s", string(data))
	}
}

func TestChange_Marshal(t *testing.T) {
	c1 := loadChange(t, "testdata/changeset_38162206.osc")
	cleanXMLNameFromChange(c1)
	data, err := c1.Marshal()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	c2, err := UnmarshalChange(data)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if !reflect.DeepEqual(c1, c2) {
		t.Errorf("changes are not equal")
		t.Logf("%+v", c1)
		t.Logf("%+v", c2)
	}

	// second changeset
	c1 = loadChange(t, "testdata/changeset_38162210.osc")
	cleanXMLNameFromChange(c1)
	data, err = c1.Marshal()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	c2, err = UnmarshalChange(data)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if !reflect.DeepEqual(c1, c2) {
		t.Errorf("changes are not equal")
		t.Logf("%+v", c1)
		t.Logf("%+v", c2)
	}

	// minute diff change
	c1 = loadChange(t, "testdata/minute_871.osc")
	cleanXMLNameFromChange(c1)
	data, err = c1.Marshal()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	c2, err = UnmarshalChange(data)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if !reflect.DeepEqual(c1, c2) {
		t.Errorf("changes are not equal")
		t.Logf("%+v", c1)
		t.Logf("%+v", c2)
	}

	// empty change
	c3 := &Change{}
	data, err = c3.Marshal()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	if l := len(data); l != 0 {
		t.Errorf("empty should be empty, got %v", l)
	}
}

func TestChange_ToHistoryDatasource(t *testing.T) {
	ctx := context.Background()
	c := &Change{
		Create: &OSM{
			Nodes: Nodes{{ID: 1, Version: 1}},
		},
		Modify: &OSM{
			Nodes: Nodes{{ID: 2, Version: 2}},
		},
		Delete: &OSM{
			Nodes: Nodes{{ID: 3, Version: 3}},
		},
	}
	ds := c.ToHistoryDatasource()

	n1, err := ds.NodeHistory(ctx, 1)
	if err != nil {
		t.Fatalf("get error: %v", err)
	}

	if !n1[0].Visible {
		t.Errorf("created node should be visible")
	}

	n2, err := ds.NodeHistory(ctx, 2)
	if err != nil {
		t.Fatalf("get error: %v", err)
	}

	if !n2[0].Visible {
		t.Errorf("modified node should be visible")
	}

	n3, err := ds.NodeHistory(ctx, 3)
	if err != nil {
		t.Fatalf("get error: %v", err)
	}

	if n3[0].Visible {
		t.Errorf("deleted node should not be visible")
	}
}

func cleanXMLNameFromChange(c *Change) {
	c.Version = 0
	c.Generator = ""
	c.Copyright = ""
	c.Attribution = ""
	c.License = ""
	if c.Create != nil {
		cleanXMLNameFromOSM(c.Create)
	}
	if c.Modify != nil {
		cleanXMLNameFromOSM(c.Modify)
	}
	if c.Delete != nil {
		cleanXMLNameFromOSM(c.Delete)
	}
}
