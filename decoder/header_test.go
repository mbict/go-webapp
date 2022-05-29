package decoder

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type headerChild struct {
	String string `header:"string"`
}

type headerTest struct {
	String    string   `header:"string"`
	StringPtr *string  `header:"string"`
	Int       int      `header:"int"`
	Int8      int8     `header:"int8"`
	Int16     int16    `header:"int16"`
	Int32     int32    `header:"int32"`
	Int64     int64    `header:"int64"`
	Uint      uint     `header:"uint"`
	Uint8     uint8    `header:"uint8"`
	Uint16    uint16   `header:"uint16"`
	Uint32    uint32   `header:"uint32"`
	Uint64    uint64   `header:"uint64"`
	Float32   float32  `header:"float32"`
	Float64   float64  `header:"float64"`
	Bool      bool     `header:"bool"`
	Strings   []string `header:"strings"`
	Nested    headerChild
	NestedPtr *headerChild
}

func TestHeaderDecoder(t *testing.T) {
	dec, err := NewHeaderDecoder(headerTest{}, "header")

	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "/foo", nil)
	req.Header.Add("string", "stringVal")
	req.Header.Add("int", "1")
	req.Header.Add("int8", "8")
	req.Header.Add("int16", "16")
	req.Header.Add("int32", "32")
	req.Header.Add("int64", "64")
	req.Header.Add("uint", "10")
	req.Header.Add("uint8", "80")
	req.Header.Add("uint16", "160")
	req.Header.Add("uint32", "320")
	req.Header.Add("uint64", "640")
	req.Header.Add("float32", "12.34")
	req.Header.Add("float64", "45.67")
	req.Header.Add("bool", "true")

	assert.NoError(t, err)

	out := &headerTest{}
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
