package decoder

import (
	"context"
	"github.com/mbict/httprouter"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type pathChild struct {
	String string `path:"string"`
}

type pathTest struct {
	String    string   `path:"string"`
	StringPtr *string  `path:"string"`
	Int       int      `path:"int"`
	Int8      int8     `path:"int8"`
	Int16     int16    `path:"int16"`
	Int32     int32    `path:"int32"`
	Int64     int64    `path:"int64"`
	Uint      uint     `path:"uint"`
	Uint8     uint8    `path:"uint8"`
	Uint16    uint16   `path:"uint16"`
	Uint32    uint32   `path:"uint32"`
	Uint64    uint64   `path:"uint64"`
	Float32   float32  `path:"float32"`
	Float64   float64  `path:"float64"`
	Bool      bool     `path:"bool"`
	Strings   []string `path:"strings"`
	Nested    pathChild
	NestedPtr *pathChild
}

func TestPathDecoder(t *testing.T) {
	dec, err := NewPathDecoder(pathTest{}, "path")

	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "/foo", nil)
	params := httprouter.Params{
		httprouter.Param{Key: "string", Value: "stringVal"},
		httprouter.Param{Key: "int", Value: "1"},
		httprouter.Param{Key: "int8", Value: "8"},
		httprouter.Param{Key: "int16", Value: "16"},
		httprouter.Param{Key: "int32", Value: "32"},
		httprouter.Param{Key: "int64", Value: "64"},
		httprouter.Param{Key: "uint", Value: "10"},
		httprouter.Param{Key: "uint8", Value: "80"},
		httprouter.Param{Key: "uint16", Value: "160"},
		httprouter.Param{Key: "uint32", Value: "320"},
		httprouter.Param{Key: "uint64", Value: "640"},
		httprouter.Param{Key: "float32", Value: "12.34"},
		httprouter.Param{Key: "float64", Value: "45.67"},
		httprouter.Param{Key: "bool", Value: "true"},
	}
	req = req.WithContext(context.WithValue(context.Background(), httprouter.ParamsKey, params))

	assert.NoError(t, err)

	out := &pathTest{}
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
