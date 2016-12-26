package logberry

import (
	"bytes"
	"testing"
	"github.com/stretchr/testify/require"
	"encoding/json"
)

func Test_EventData_Basic(t *testing.T) {
	require := require.New(t)
	
	data := EventDataMap{
		"Alice": EventDataString("Miranda"),
		"Kobayashi": EventDataString("Maru"),
		"Star Trek": EventDataMap{
			"Season": EventDataInt32(5),
			"Episode #": EventDataInt32(2),
			"Title": EventDataString("Darmok"),
			"Quality": EventDataFloat32(9.9),
		},
	}

	buff := new(bytes.Buffer)
	data.WriteRecurse(buff)
	t.Log(buff.String())

	test := "{ Alice=\"Miranda\" Kobayashi=\"Maru\" \"Star Trek\"={ Season=5 \"Episode #\"=2 Title=\"Darmok\" Quality=9.9 } }"

	require.Equal(buff.String(), test)
	
}

func Test_EventData_JSON(t *testing.T) {
	require := require.New(t)
	
	data := EventDataMap{
		"Alice": EventDataString("Miranda"),
		"Kobayashi": EventDataString("Maru"),
		"Star Trek": EventDataMap{
			"Season": EventDataInt32(5),
			"Episode #": EventDataInt32(2),
			"Title": EventDataString("Darmok"),
			"Quality": EventDataFloat32(9.9),
		},
	}

	b, err := json.Marshal(data)
	require.Nil(err)

	t.Log(string(b))
	require.Equal(string(b), "{\"Alice\":\"Miranda\",\"Kobayashi\":\"Maru\",\"Star Trek\":{\"Episode #\":2,\"Quality\":9.9,\"Season\":5,\"Title\":\"Darmok\"}}")

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
	data.WriteRecurse(buff)
	t.Log(buff.String())

	test := "{ Name=\"Alice Miranda\" Friends=[ \"Anna\", \"Black Bear\", { First=\"Melvin\" Last=\"Hedgehog\" }, \"Miles\" ] }"

	require.Equal(buff.String(), test)
	
}
