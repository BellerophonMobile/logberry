package logberry

import (
	"bytes"
	"encoding/json"
	"testing"
)

type expectation struct {
	ex string
	v  *D
}

//----------------------------------------------------------------------
//----------------------------------------------------------------------
func runcases(tests []expectation, t *testing.T) {

	for _, c := range tests {
		str, err := json.Marshal(c.v)
		if err != nil {
			t.Error("Unexpected error", err)
		}
		if bytes.Compare(str, []byte(c.ex)) != 0 {
			t.Errorf("Expected '%v', got '%v'", c.ex, string(str))
		}
	}

	// end runcases
}

//----------------------------------------------------------------------
//----------------------------------------------------------------------
type test struct {
	StringField string
	IntField    int
}

func TestDBuild(t *testing.T) {

	var tests = []expectation{

		{
			v:  DBuild(nil),
			ex: "{}",
		},

		{
			v:  DBuild(8),
			ex: "{\"value\":8}",
		},

		{
			v:  DBuild("mushi"),
			ex: "{\"value\":\"mushi\"}",
		},

		{
			v:  DBuild(&test{StringField: "Banana", IntField: 7}),
			ex: "{\"IntField\":7,\"StringField\":\"Banana\"}",
		},

		{
			v:  DBuild(&D{"Fruit": "Banana"}),
			ex: "{\"Fruit\":\"Banana\"}",
		},

		{
			v:  DBuild(D{"Fruit": "Banana"}),
			ex: "{\"Fruit\":\"Banana\"}",
		},

		{
			v:  DBuild(map[string]int{"Sector": 12, "System": 4}),
			ex: "{\"Sector\":12,\"System\":4}",
		},

		{
			v:  DBuild(map[int]string{12: "Joe", 4: "Tom"}),
			ex: "{\"12\":\"Joe\",\"4\":\"Tom\"}",
		},
	}

	runcases(tests, t)

	// end TestDBuild
}

//----------------------------------------------------------------------
func TestCopyFromD(t *testing.T) {

	var tests = []expectation{

		{
			v:  (&D{}).CopyFromD(nil),
			ex: "{}",
		},

		{
			v:  (&D{}).CopyFromD(&D{"Field1": "Data"}),
			ex: "{\"Field1\":\"Data\"}",
		},

		{
			v:  (&D{"Field2": "Atad"}).CopyFromD(&D{"Field1": "Data"}),
			ex: "{\"Field1\":\"Data\",\"Field2\":\"Atad\"}",
		},

		{
			v:  (&D{"Field2": "Atad", "Field3": "Foo"}).CopyFromD(&D{"Field1": "Data"}),
			ex: "{\"Field1\":\"Data\",\"Field2\":\"Atad\",\"Field3\":\"Foo\"}",
		},

		{
			v:  (&D{"Field2": "Atad", "Field3": "Foo"}).CopyFromD(&D{"Field1": "Data", "Field4": "Bar"}),
			ex: "{\"Field1\":\"Data\",\"Field2\":\"Atad\",\"Field3\":\"Foo\",\"Field4\":\"Bar\"}",
		},
	}

	runcases(tests, t)

	// TestCopyFromD
}

//----------------------------------------------------------------------
func TestCopyFrom(t *testing.T) {

	var tests = []expectation{

		{
			v:  (&D{}).CopyFrom(nil),
			ex: "{}",
		},

		{
			v:  (&D{"Field2": "Atad"}).CopyFrom(8),
			ex: "{\"Field2\":\"Atad\",\"value\":8}",
		},

		{
			v:  (&D{"Field2": "Atad"}).CopyFrom(8).CopyFrom("mushi"),
			ex: "{\"Field2\":\"Atad\",\"value\":[8,\"mushi\"]}",
		},

		{
			v:  (&D{"Field2": "Atad", "Field1": "Data"}).CopyFrom(&test{StringField: "Banana", IntField: 7}),
			ex: "{\"Field1\":\"Data\",\"Field2\":\"Atad\",\"IntField\":7,\"StringField\":\"Banana\"}",
		},
	}

	runcases(tests, t)

	// end TestCopyFrom
}

//----------------------------------------------------------------------
func TestAggregateFrom(t *testing.T) {

	var tests = []expectation{

		{
			v:  (&D{}).AggregateFrom([]interface{}{nil}),
			ex: "{}",
		},

		{
			v:  (&D{}).AggregateFrom([]interface{}{8, 9, "mushi"}),
			ex: "{\"value\":[8,9,\"mushi\"]}",
		},

		{
			v:  (&D{}).AggregateFrom([]interface{}{8, &D{"Fruit": "Banana"}, 9, "mushi"}),
			ex: "{\"Fruit\":\"Banana\",\"value\":[8,9,\"mushi\"]}",
		},

		{
			v:  (&D{}).AggregateFrom([]interface{}{8, &D{"Fruit": "Banana"}, 9, &D{"Fruit": "Candy"}}),
			ex: "{\"Fruit\":\"Candy\",\"value\":[8,9]}",
		},

		// Above are (or should be) same tests as TestDAggregate

		{
			v:  (&D{"Ship": "Black Pearl"}).AggregateFrom([]interface{}{nil}),
			ex: "{\"Ship\":\"Black Pearl\"}",
		},

		{
			v:  (&D{"Ship": "Black Pearl"}).AggregateFrom([]interface{}{8, 9, "mushi"}),
			ex: "{\"Ship\":\"Black Pearl\",\"value\":[8,9,\"mushi\"]}",
		},
	}

	runcases(tests, t)

	// end TestAggregateFrom
}

//----------------------------------------------------------------------
func TestDAggregate(t *testing.T) {

	var tests = []expectation{

		{
			v:  DAggregate([]interface{}{nil}),
			ex: "{}",
		},

		{
			v:  DAggregate([]interface{}{8, 9, "mushi"}),
			ex: "{\"value\":[8,9,\"mushi\"]}",
		},

		{
			v:  DAggregate([]interface{}{8, &D{"Fruit": "Banana"}, 9, "mushi"}),
			ex: "{\"Fruit\":\"Banana\",\"value\":[8,9,\"mushi\"]}",
		},

		{
			v:  DAggregate([]interface{}{8, &D{"Fruit": "Banana"}, 9, &D{"Fruit": "Candy"}}),
			ex: "{\"Fruit\":\"Candy\",\"value\":[8,9]}",
		},
	}

	runcases(tests, t)

	// end TestDAggregate
}
