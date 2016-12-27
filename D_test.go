package logberry

import (
	"bytes"
	"testing"
)

type d_expectation struct {
	ex string
	v  interface{}
}

func runcases_d(tests []d_expectation, t *testing.T) {

	for _, c := range tests {
		t.Run(c.ex, func(t *testing.T) {
			eventdata := Copy(c.v)

			buff := new(bytes.Buffer)
			eventdata.WriteTo(buff)

			if bytes.Compare(buff.Bytes(), []byte(c.ex)) != 0 {
				t.Errorf("Expected '%v', got '%v'", c.ex, buff.String())
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
			v:  nil,
			ex: "{ }",
		},

		{
			v:  mushi1,
			ex: "{ }",
		},

		{
			v:  8,
			ex: "8",
		},

		{
			v:  "mushi",
			ex: "\"mushi\"",
		},

		{
			v:  &test{StringField: "Banana", IntField: 7},
			ex: "{ IntField=7 StringField=\"Banana\" }",
		},

		{
			v:  D{"Fruit": "Banana"},
			ex: "{ Fruit=\"Banana\" }",
		},

		{
			v:  map[string]int{"Sector": 12, "System": 4},
			ex: "{ Sector=12 System=4 }",
		},

		{
			v:  map[int]string{12: "Joe", 4: "Tom"},
			ex: "{ 12=\"Joe\" 4=\"Tom\" }",
		},

		{
			v:  &test2{"mushi", 4},
			ex: "{ PublicField=4 }",
		},

		{
			v:  D{"Field1": "Data"},
			ex: "{ Field1=\"Data\" }",
		},

		{
			v:  D{"Field2": "Atad", "Field1": "Data"},
			ex: "{ Field1=\"Data\" Field2=\"Atad\" }",
		},

		{
			v:  D{"Field2": "Atad", "Field3": "Foo", "Field1": "Data"},
			ex: "{ Field1=\"Data\" Field2=\"Atad\" Field3=\"Foo\" }",
		},

		{
			v:  []string{"foo", "bar"},
			ex: "[\"foo\", \"bar\"]",
		},

		{
			v:  D{"baz": []string{"foo", "bar"}},
			ex: "{ baz=[\"foo\", \"bar\"] }",
		},
	}

	runcases_d(tests, t)

}
