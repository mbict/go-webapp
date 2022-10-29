package decoder

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type cookieChild struct {
	String string `cookie:"string"`
}

type cookieTest struct {
	String    string   `cookie:"string"`
	StringPtr *string  `cookie:"string"`
	Int       int      `cookie:"int"`
	Int8      int8     `cookie:"int8"`
	Int16     int16    `cookie:"int16"`
	Int32     int32    `cookie:"int32"`
	Int64     int64    `cookie:"int64"`
	Uint      uint     `cookie:"uint"`
	Uint8     uint8    `cookie:"uint8"`
	Uint16    uint16   `cookie:"uint16"`
	Uint32    uint32   `cookie:"uint32"`
	Uint64    uint64   `cookie:"uint64"`
	Float32   float32  `cookie:"float32"`
	Float64   float64  `cookie:"float64"`
	Bool      bool     `cookie:"bool"`
	Strings   []string `cookie:"strings"`
	Nested    cookieChild
	NestedPtr *cookieChild
}

func TestCookieDecoder(t *testing.T) {
	dec, err := NewCookieDecoder(cookieTest{}, "cookie")

	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "/foo", nil)
	req.AddCookie(&http.Cookie{Name: "string", Value: "stringVal"})
	req.AddCookie(&http.Cookie{Name: "int", Value: "1"})
	req.AddCookie(&http.Cookie{Name: "int8", Value: "8"})
	req.AddCookie(&http.Cookie{Name: "int16", Value: "16"})
	req.AddCookie(&http.Cookie{Name: "int32", Value: "32"})
	req.AddCookie(&http.Cookie{Name: "int64", Value: "64"})
	req.AddCookie(&http.Cookie{Name: "uint", Value: "10"})
	req.AddCookie(&http.Cookie{Name: "uint8", Value: "80"})
	req.AddCookie(&http.Cookie{Name: "uint16", Value: "160"})
	req.AddCookie(&http.Cookie{Name: "uint32", Value: "320"})
	req.AddCookie(&http.Cookie{Name: "uint64", Value: "640"})
	req.AddCookie(&http.Cookie{Name: "float32", Value: "12.34"})
	req.AddCookie(&http.Cookie{Name: "float64", Value: "45.67"})
	req.AddCookie(&http.Cookie{Name: "bool", Value: "true"})

	assert.NoError(t, err)

	out := &cookieTest{}
	err = dec(req, out)

	assert.NoError(t, err)

	assert.Equal(t, "stringVal", out.String)
	assert.Equal(t, asPtr("stringVal"), out.StringPtr)
	assert.Equal(t, 1, out.Int)
	assert.Equal(t, int8(8), out.Int8)
	assert.Equal(t, int16(16), out.Int16)
	assert.Equal(t, int32(32), out.Int32)
	assert.Equal(t, int64(64), out.Int64)
	assert.Equal(t, uint(10), out.Uint)
	assert.Equal(t, uint8(80), out.Uint8)
	assert.Equal(t, uint16(160), out.Uint16)
	assert.Equal(t, uint32(320), out.Uint32)
	assert.Equal(t, uint64(640), out.Uint64)
	assert.Equal(t, float32(12.34), out.Float32)
	assert.Equal(t, float64(45.67), out.Float64)
	assert.Equal(t, true, out.Bool)
}
