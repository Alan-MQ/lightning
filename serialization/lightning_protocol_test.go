package serialization

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

type TestStruct struct {
	Int        int
	String     string
	unExported int
}

func TestMarshal(t *testing.T) {
	s := TestStruct{Int: 1, String: "something", unExported: 112}
	stream, err := Marshal(s)
	if err != nil {
		// panic(err)
		panic(err)
	}
	err = UnMarshal(&TestStruct{}, stream)
	// assert.Equal(err, nil)
	assert.Equal(t, nil, err)
}
