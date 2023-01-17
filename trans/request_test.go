package trans

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeader(t *testing.T) {
	asserts := assert.New(t)
	httpReq, err := http.NewRequest(http.MethodGet, "https://example.com/", nil)
	if !asserts.Nil(err, err) {
		return
	}
	httpReq.Header.Add("Content-Type", "application/json")
	req, err := NewCurlRequest(httpReq)
	if !asserts.Nil(err, err) {
		return
	}
	headers := req.GetHeaders()
	asserts.Equal(headers.Get("Content-Type"), "application/json", "unable to get the correct header value")
	req.AddHeader("Accept-Language", "*")
	asserts.Contains(req.GetHeader("Accept-Language"), "*", "unable to add the header")
}

func TestGetCommands(t *testing.T) {
	asserts := assert.New(t)
	httpReq, err := http.NewRequest(http.MethodGet, "https://example.com/", nil)
	if !asserts.Nil(err, err) {
		return
	}
	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Accept", "application/json")
	expected := `curl -s -i -k -H Content-Type: application/json -H Accept: application/json -X GET https://example.com/`
	req, err := NewCurlRequest(httpReq)
	if !asserts.Nil(err, err) {
		return
	}
	err = req.GenerateCommand(true)
	if !asserts.Nil(err, err) {
		return
	}
	cmd := req.GetCommands()
	asserts.Equal(expected, strings.Join(cmd, " "), "incorrect command")
}
