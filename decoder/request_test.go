package decoder

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type requestChild struct {
	String string `request:"string"`
}

type requestTest struct {
	RemoteAddr string `request:"remote-addr"`
	Host       string `request:"host"`
	Method     string `request:"method"`
	Url        string `request:"url"`
	UrlHost    string `request:"url:host"`
	UrlQuery   string `request:"url:query"`
	UrlPath    string `request:"url:path"`
	UrlSchema  string `request:"url:scheme"`
}

func TestRequestDecoder(t *testing.T) {
	dec, err := NewRequestDecoder(requestTest{}, "request")

	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "https://test.com/foo?bar=baz&foo=1", nil)

	assert.NoError(t, err)

	out := &requestTest{}
	err = dec(req, out)

	assert.NoError(t, err)

	assert.Equal(t, "", out.RemoteAddr)
	assert.Equal(t, "test.com", out.Host)
	assert.Equal(t, "GET", out.Method)
	assert.Equal(t, "https://test.com/foo?bar=baz&foo=1", out.Url)
	assert.Equal(t, "test.com", out.UrlHost)
	assert.Equal(t, "bar=baz&foo=1", out.UrlQuery)
	assert.Equal(t, "/foo", out.UrlPath)
	assert.Equal(t, "https", out.UrlSchema)

}
