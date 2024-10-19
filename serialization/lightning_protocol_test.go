package serialization

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

type testStruct struct {
	Int        int
	String     string
	unExported int
}

func TestMarshal(t *testing.T) {
	s := testStruct{Int: 1, String: "something", unExported: 112}
	stream, err := Marshal(s)
	obj, err := UnMarshal(&testStruct{}, stream))
	// assert.Equal(err, nil)
	assert.Equal(t, nil, err)
}
