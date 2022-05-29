package decoder

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type queryChild struct {
	String string `query:"string"`
}

type queryTest struct {
	String    string   `query:"string"`
	StringPtr *string  `query:"string"`
	Int       int      `query:"int"`
	Int8      int8     `query:"int8"`
	Int16     int16    `query:"int16"`
	Int32     int32    `query:"int32"`
	Int64     int64    `query:"int64"`
	Uint      uint     `query:"uint"`
	Uint8     uint8    `query:"uint8"`
	Uint16    uint16   `query:"uint16"`
	Uint32    uint32   `query:"uint32"`
	Uint64    uint64   `query:"uint64"`
	Float32   float32  `query:"float32"`
	Float64   float64  `query:"float64"`
	Bool      bool     `query:"bool"`
	Strings   []string `query:"strings"`
	Nested    queryChild
	NestedPtr *queryChild
}

func TestQueryDecoder(t *testing.T) {
	dec, err := NewQueryDecoder(queryTest{}, "query")

	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "/foo?string=stringVal&int=1&int8=8&int16=16&int32=32&int64=64&uint=10&uint8=80&uint16=160&uint32=320&uint64=640&float32=12.34&float64=45.67&bool=true", nil)

	assert.NoError(t, err)

	out := &queryTest{}
	err = dec(req, out)

	assert.NoError(t, err)

	assert.Equal(t, out.String, "stringVal")
	assert.Equal(t, out.StringPtr, asPtr("stringVal"))
	assert.Equal(t, out.Int, 1)
	assert.Equal(t, out.Int8, int8(8))
	assert.Equal(t, out.Int16, int16(16))
	assert.Equal(t, out.Int32, int32(32))
	assert.Equal(t, out.Int64, int64(64))
	assert.Equal(t, out.Uint, uint(10))
	assert.Equal(t, out.Uint8, uint8(80))
	assert.Equal(t, out.Uint16, uint16(160))
	assert.Equal(t, out.Uint32, uint32(320))
	assert.Equal(t, out.Uint64, uint64(640))
	assert.Equal(t, out.Float32, float32(12.34))
	assert.Equal(t, out.Float64, float64(45.67))
	assert.Equal(t, out.Bool, true)
}

func asPtr[T comparable](in T) *T {
	return &in
}
