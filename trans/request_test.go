package trans

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/0x726f6f6b6965/go-curl/filetype"
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
	req, err := NewCurlRequest(httpReq)
	if !asserts.Nil(err, err) {
		return
	}
	err = req.GenerateCommand(true)
	if !asserts.Nil(err, err) {
		return
	}
	cmd := strings.Join(req.GetCommands(), " ")
	asserts.Contains(cmd, "-H Content-Type: application/json", "incorrect command of header")
	asserts.Contains(cmd, "-H Accept: application/json", "incorrect command of header")
	asserts.Contains(cmd, "-X GET https://example.com/", "incorrect command of method and url")
	asserts.Contains(cmd, "curl -s -i -k ", "incorrect command")
}

func TestFilename(t *testing.T) {
	asserts := assert.New(t)
	output := &testBody{User: "andy", Flag: true, Count: 9}
	body, err := json.Marshal(output)
	if !asserts.Nil(err, err) {
		return
	}
	httpReq, err := http.NewRequest(http.MethodPost, "https://example.com/", bytes.NewReader(body))
	if !asserts.Nil(err, err) {
		return
	}
	httpReq.Header.Add("Content-Type", "application/json")
	req, err := NewCurlRequest(httpReq)
	if !asserts.Nil(err, err) {
		return
	}
	err = req.GenerateCommand(true)
	if !asserts.Nil(err, err) {
		return
	}
	filename := req.GetFilename()
	data, err := os.ReadFile(filename)
	if !asserts.Nil(err, err) {
		return
	}
	defer os.Remove(filename)
	var fileBody testBody
	err = json.Unmarshal(data, &fileBody)
	if !asserts.Nil(err, err) {
		return
	}
	asserts.EqualValues(output, &fileBody)
}

func TestSetHeaderAfterGeneratedCmd(t *testing.T) {
	asserts := assert.New(t)
	output := &testBody{User: "andy", Flag: true, Count: 9}
	body, err := json.Marshal(output)
	if !asserts.Nil(err, err) {
		return
	}
	httpReq, err := http.NewRequest(http.MethodPost, "https://example.com/", bytes.NewReader(body))
	if !asserts.Nil(err, err) {
		return
	}
	httpReq.Header.Add("Content-Type", "application/json")
	req, err := NewCurlRequest(httpReq)
	if !asserts.Nil(err, err) {
		return
	}
	err = req.GenerateCommand(true)
	if !asserts.Nil(err, err) {
		return
	}
	err = req.AddHeader("Accept", "application/json")
	if !asserts.Nil(err, err) {
		return
	}
	os.Remove(req.GetFilename())
	asserts.Contains(strings.Join(req.GetCommands(), " "), "-H Accept: application/json", "incorrect command of header")
}

func TestSetPrivatekey(t *testing.T) {
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
	req.SetPrivateKey(filetype.PEM, "./abc.PEM")

	err = req.GenerateCommand(true)
	if !asserts.Nil(err, err) {
		return
	}
	cmd := strings.Join(req.GetCommands(), " ")
	asserts.Contains(cmd, "--key-type PEM", "incorrect command of private key type")
	asserts.Contains(cmd, "--key ./abc.PEM", "incorrect command of private key path")
}

func TestSetCertificate(t *testing.T) {
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
	req.SetCertificate(filetype.PEM, "./abc.PEM", "pwd")

	err = req.GenerateCommand(true)
	if !asserts.Nil(err, err) {
		return
	}
	cmd := strings.Join(req.GetCommands(), " ")
	asserts.Contains(cmd, "--cert-type PEM", "incorrect command of certificate type")
	asserts.Contains(cmd, "--cert ./abc.PEM:pwd", "incorrect command of certificate path")
}

func TestSetCertificateWithoutPassword(t *testing.T) {
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
	req.SetCertificate(filetype.PEM, "./abc.PEM", "")

	err = req.GenerateCommand(true)
	if !asserts.Nil(err, err) {
		return
	}
	cmd := strings.Join(req.GetCommands(), " ")
	asserts.Contains(cmd, "--cert-type PEM", "incorrect command of certificate type")
	asserts.Contains(cmd, "--cert ./abc.PEM", "incorrect command of certificate path")
}

func TestSetCA(t *testing.T) {
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
	req.SetCA("./abc.PEM")

	err = req.GenerateCommand(true)
	if !asserts.Nil(err, err) {
		return
	}
	cmd := strings.Join(req.GetCommands(), " ")
	asserts.Contains(cmd, "--cacert ./abc.PEM", "incorrect command of ca certificate path")
}

func TestDo(t *testing.T) {
	asserts := assert.New(t)
	httpReq, err := http.NewRequest(http.MethodGet, "https://example.com/", nil)
	if !asserts.Nil(err, err) {
		return
	}
	httpReq.Header.Add("Content-Type", "application/json")
	req := &curlRequest{ctx: context.Background(), request: httpReq, cleanup: func() error { return nil }}
	output := func(ctx context.Context, name string, args ...string) ([]byte, error) {
		resp := fmt.Sprintf(`%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n\r\n%s`,
			"HTTP/1.1 200 OK", "Connection: close", "X-Frame-Options: SAMEORIGIN",
			"Cache-Control: no-cache, no-store, must-revalidate",
			"Content-Length: 123", "Content-Type: application/json",
			"{\"user\":\"response\",\"flag\":false,\"count\":5}",
		)
		return []byte(resp), nil
	}
	req.combinedOutput = output
	err = req.GenerateCommand(true)
	if !asserts.Nil(err, err) {
		return
	}
	resp, err := req.Do()
	if !asserts.Nil(err, err) {
		return
	}

	if !asserts.Equal(resp.StatusCode, http.StatusOK, "incorrect status code") {
		return
	}

	if !asserts.Equal(resp.Status, http.StatusText(http.StatusOK), "incorrect status code") {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if !asserts.Nil(err, "read body error") {
		return
	}

	var result testBody
	err = json.Unmarshal(body, &result)
	if !asserts.Nil(err, "parse body error") {
		return
	}

	if !asserts.True(result.Count == 5 && result.Flag == false && result.User == "response", "incorrect body") {
		return
	}
}

func TestExecBeforeGenerateCommand(t *testing.T) {
	asserts := assert.New(t)
	httpReq, err := http.NewRequest(http.MethodGet, "https://example.com/", nil)
	if !asserts.Nil(err, err) {
		return
	}
	httpReq.Header.Add("Content-Type", "application/json")
	req := &curlRequest{ctx: context.Background(), request: httpReq, cleanup: func() error { return nil }}
	output := func(ctx context.Context, name string, args ...string) ([]byte, error) {
		resp := fmt.Sprintf("eee")
		return []byte(resp), nil
	}
	req.combinedOutput = output
	_, err = req.Execute()
	asserts.ErrorContains(err, "unable to run this command, command")
}
