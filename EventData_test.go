package logberry

import (
	"bytes"
	"testing"
	"github.com/stretchr/testify/require"
	"encoding/json"
)

func Test_EventData_NilMap(t *testing.T) {
	require := require.New(t)

	data := EventDataMap(nil)

	buff := new(bytes.Buffer)
	data.WriteTo(buff)

	test := "{ }"

	require.Equal(buff.String(), test)
	
}

func Test_EventData_Basic(t *testing.T) {
	require := require.New(t)
	
	data := EventDataMap{
		"Alice": EventDataString("Miranda"),
		"Kobayashi": EventDataString("Maru"),
		"Star Trek": EventDataMap{
			"Season": EventDataInt64(5),
			"Episode #": EventDataInt64(2),
			"Title": EventDataString("Darmok"),
			"Quality": EventDataFloat64(9.9),
		},
	}

	buff := new(bytes.Buffer)
	data.WriteTo(buff)

	test := "{ Alice=\"Miranda\" Kobayashi=\"Maru\" \"Star Trek\"={ \"Episode #\"=2 Quality=9.9 Season=5 Title=\"Darmok\" } }"

	require.Equal(test, buff.String())
	
}

func Test_EventData_JSON(t *testing.T) {
	require := require.New(t)
	
	data := EventDataMap{
		"Alice": EventDataString("Miranda"),
		"Kobayashi": EventDataString("Maru"),
		"Star Trek": EventDataMap{
			"Season": EventDataInt64(5),
			"Episode #": EventDataInt64(2),
			"Title": EventDataString("Darmok"),
			"Quality": EventDataFloat64(9.9),
		},
	}

	b, err := json.Marshal(data)
	require.Nil(err)

	require.Equal("{\"Alice\":\"Miranda\",\"Kobayashi\":\"Maru\",\"Star Trek\":{\"Episode #\":2,\"Quality\":9.9,\"Season\":5,\"Title\":\"Darmok\"}}", string(b))

}

func Test_EventData_Slice(t *testing.T) {
	require := require.New(t)
	
	data := EventDataMap{
		"Name": EventDataString("Alice Miranda"),
		"Friends": EventDataSlice{
			EventDataString("Anna"),
			EventDataString("Black Bear"),
			EventDataMap{
				"First": EventDataString("Melvin"),
				"Last": EventDataString("Hedgehog"),
			},
			EventDataString("Miles"),
		},
	}

	buff := new(bytes.Buffer)
	data.WriteTo(buff)

	test := "{ Friends=[ \"Anna\", \"Black Bear\", { First=\"Melvin\" Last=\"Hedgehog\" }, \"Miles\" ] Name=\"Alice Miranda\" }"

	require.Equal(test, buff.String())
	
}
