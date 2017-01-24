package logberry

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func Test_Failure(t *testing.T) {
	require := require.New(t)

	key := "status"
	
	err := Main.Failure("Epic societal collapse", D{key: 404})

	t.Log("D", err.Data)
	d, ok := err.Data[key]
	require.True(ok)
	
	num, ok := d.(EventDataInt64)
	require.True(ok)

	require.Equal(num, EventDataInt64(404))
	
}
