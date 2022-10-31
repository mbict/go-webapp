package decoder

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
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

type queryExplodeTest struct {
	StringSlice    []string `query:"strings"`
	CommaSlice     []string `query:"stringsc,delimiter:comma"`
	SemiColonSlice []string `query:"stringssc,delimiter:semicolon"`
	SpaceSlice     []string `query:"stringss,delimiter:space"`
	PipeSlice      []string `query:"stringsp,delimiter:pipe"`
}

func TestQueryDecoderExplodeSlice(t *testing.T) {
	dec, err := NewQueryDecoder(queryExplodeTest{}, "query")

	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "/foo?strings=abc&strings=foo&stringsc=foo,bar,baz&stringss=foo%20bar%20baz&stringssc=foo%3Bbar%3Bbaz&stringsp=foo|bar|baz", nil)

	assert.NoError(t, err)

	out := &queryExplodeTest{}
	err = dec(req, out)

	assert.NoError(t, err)

	assert.Len(t, out.StringSlice, 2)
	assert.Equal(t, []string{"abc", "foo"}, out.StringSlice)
	assert.Equal(t, []string{"foo", "bar", "baz"}, out.CommaSlice)
	assert.Equal(t, []string{"foo", "bar", "baz"}, out.SemiColonSlice)
	assert.Equal(t, []string{"foo", "bar", "baz"}, out.SpaceSlice)
	assert.Equal(t, []string{"foo", "bar", "baz"}, out.PipeSlice)
}

type querySlicedTextMarshallers struct {
	UUIDS    []uuid.UUID `query:"ids"`
	UUIDSSTR []uuid.UUID `query:"idsc,delimiter:comma"`
}

func TestQueryDecoderSlicedTextMarshallers(t *testing.T) {
	dec, err := NewQueryDecoder(querySlicedTextMarshallers{}, "query")

	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "/foo?ids=ab75db78-c9eb-4582-bd2c-02a7ed53c80d&ids=985ac730-7ebc-4099-a463-c8a59fb116a2&idsc=ab75db78-c9eb-4582-bd2c-02a7ed53c80d,985ac730-7ebc-4099-a463-c8a59fb116a2", nil)

	assert.NoError(t, err)

	out := &querySlicedTextMarshallers{}
	err = dec(req, out)

	assert.NoError(t, err)

	expectedUUIDs := []uuid.UUID{
		uuid.MustParse("ab75db78-c9eb-4582-bd2c-02a7ed53c80d"),
		uuid.MustParse("985ac730-7ebc-4099-a463-c8a59fb116a2"),
	}

	assert.Equal(t, expectedUUIDs, out.UUIDS)
	assert.Equal(t, expectedUUIDs, out.UUIDSSTR)
}

type queryDataTime struct {
	Date  time.Time   `query:"datetime"`
	Dates []time.Time `query:"datetimes,delimiter:comma"`
}

func TestQueryDecoderDates(t *testing.T) {
	dec, err := NewQueryDecoder(queryDataTime{}, "query")

	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "/foo?datetime=2006-01-02T15:04:05.000Z&datetimes=2006-01-02T15:04:05.000Z,2016-01-02T15:04:05.000Z", nil)

	assert.NoError(t, err)

	out := &queryDataTime{}
	err = dec(req, out)

	assert.NoError(t, err)

	firstDate, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05.000Z")
	secondDate, _ := time.Parse(time.RFC3339, "2016-01-02T15:04:05.000Z")

	assert.Equal(t, firstDate, out.Date)
	assert.Equal(t, []time.Time{firstDate, secondDate}, out.Dates)

}
