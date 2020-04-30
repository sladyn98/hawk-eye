package github

import (
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	okResponse = `{
		"total_count":2,
		"check_runs":[{"id":570295330,"status":"completed","conclusion":"success","app":{"id":67,"name":"Travis CI"}},
					  {"id":570295331,"status":"completed","conclusion":"success","app":{"id":67,"name":"Github Actions"}}
		]
	 }`
)

func githubMetaDataResponseStub() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(okResponse))
	}))
}

func TestGetCIStatus(t *testing.T) {
	server := githubMetaDataResponseStub()
	defer server.Close()
	got, err := getCIStatus(server.URL, "hawk-eye", "sladyn98", "", "f0496ae48d0f21ccc0ef23502ccea96dd68c7938", "")
	if err != nil {
		t.Error("Something is correct", err)
	}
	assert.Equal(t, got, "true")
}
