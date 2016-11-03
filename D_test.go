package logberry

import (
	"bytes"
	"encoding/json"
	"testing"
)

type d_expectation struct {
	ex string
	v  D
}

func runcases_d(tests []d_expectation, t *testing.T) {

	for _, c := range tests {
		t.Run(c.ex, func(t *testing.T) {
			str, err := json.Marshal(c.v)
			if err != nil {
				t.Error("Unexpected error", err)
			}
			if bytes.Compare(str, []byte(c.ex)) != 0 {
				t.Errorf("Expected '%v', got '%v'", c.ex, string(str))
			}
		})
	}

}

type test struct {
	StringField string
	IntField    int
}

type test2 struct {
	privatefield string
	PublicField  int
}

func TestCopyFrom(t *testing.T) {

	// Data used in tests
	var mushi1 *test

	var tests = []d_expectation{

		{
			v:  (D{}).CopyFrom(nil),
			ex: "{}",
		},

		{
			v:  (D{}).CopyFrom(mushi1),
			ex: "{}",
		},

		{
			v:  (D{}).CopyFrom(8),
			ex: "{\"value\":8}",
		},

		{
			v:  (D{}).CopyFrom("mushi"),
			ex: "{\"value\":\"mushi\"}",
		},

		{
			v:  (D{}).CopyFrom(&test{StringField: "Banana", IntField: 7}),
			ex: "{\"IntField\":7,\"StringField\":\"Banana\"}",
		},

		{
			v:  (D{}).CopyFrom(D{"Fruit": "Banana"}),
			ex: "{\"Fruit\":\"Banana\"}",
		},

		{
			v:  (D{}).CopyFrom(map[string]int{"Sector": 12, "System": 4}),
			ex: "{\"Sector\":12,\"System\":4}",
		},

		{
			v:  (D{}).CopyFrom(map[int]string{12: "Joe", 4: "Tom"}),
			ex: "{\"12\":\"Joe\",\"4\":\"Tom\"}",
		},

		{
			v:  (D{}).CopyFrom(&test2{"mushi", 4}),
			ex: "{\"PublicField\":4}",
		},

		{
			v:  (D{}).CopyFrom(D{"Field1": "Data"}),
			ex: "{\"Field1\":\"Data\"}",
		},

		{
			v:  (D{"Field2": "Atad"}).CopyFrom(D{"Field1": "Data"}),
			ex: "{\"Field1\":\"Data\",\"Field2\":\"Atad\"}",
		},

		{
			v:  (D{"Field2": "Atad", "Field3": "Foo"}).CopyFrom(D{"Field1": "Data"}),
			ex: "{\"Field1\":\"Data\",\"Field2\":\"Atad\",\"Field3\":\"Foo\"}",
		},

		{
			v:  (D{"Field2": "Atad", "Field3": "Foo"}).CopyFrom(D{"Field1": "Data", "Field4": "Bar"}),
			ex: "{\"Field1\":\"Data\",\"Field2\":\"Atad\",\"Field3\":\"Foo\",\"Field4\":\"Bar\"}",
		},
	}

	runcases_d(tests, t)

}

func TestDAggregate(t *testing.T) {

	var tests = []d_expectation{

		{
			v:  DAggregate([]interface{}{nil}),
			ex: "{}",
		},

		{
			v:  DAggregate([]interface{}{8, 9, "mushi"}),
			ex: "{\"value\":[8,9,\"mushi\"]}",
		},

		{
			v:  DAggregate([]interface{}{8, D{"Fruit": "Banana"}, 9, "mushi"}),
			ex: "{\"Fruit\":\"Banana\",\"value\":[8,9,\"mushi\"]}",
		},

		{
			v:  DAggregate([]interface{}{8, D{"Fruit": "Banana"}, 9, D{"Fruit": "Candy"}}),
			ex: "{\"Fruit\":\"Candy\",\"value\":[8,9]}",
		},

		{
			v:  DAggregate([]interface{}{D{"Fruit": "Banana"}, 8, 9, D{"LP": "Help"}}),
			ex: "{\"Fruit\":\"Banana\",\"LP\":\"Help\",\"value\":[8,9]}",
		},
	}

	runcases_d(tests, t)

}
