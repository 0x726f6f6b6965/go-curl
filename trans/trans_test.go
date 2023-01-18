package trans

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCommandWithBasicAuth(t *testing.T) {
	asserts := assert.New(t)
	req, _ := http.NewRequest(http.MethodGet, "https://localhost:8080/home", nil)
	user, pwd := "admin", "admin123"
	req.SetBasicAuth(user, pwd)
	req.BasicAuth()
	cmd, _, _, err := GetCommand(req, true)
	if !asserts.Nil(err, err) {
		return
	}
	for i, val := range cmd {
		if strings.Contains(val, "Authorization") {
			if !asserts.Equal(cmd[i-1], "-H", "incorrect header format") {
				return
			}

			auths := strings.Split(val, ":")
			if !asserts.Equal(len(auths), 2, "incorrect header value") {
				return
			}

			if !asserts.True(strings.HasPrefix(auths[1], " Basic "), "incorrect prefix of the header value") {
				return
			}

			c, err := base64.StdEncoding.DecodeString(strings.Split(auths[1], " ")[2])
			if !asserts.Nil(err, "unable decode the header value") {
				return
			}

			cs := string(c)
			username, password, ok := strings.Cut(cs, ":")
			if !asserts.True(ok, "unable parse the auth") {
				return
			}
			if !asserts.True(username == user && password == pwd, "incorrect username/password") {
				return
			}
			return
		}
	}
	asserts.Fail("unable get the auth")
}

func TestGetCommandWithBody(t *testing.T) {
	asserts := assert.New(t)
	body := &testBody{
		User:  "admin",
		Flag:  true,
		Count: 3,
	}
	jsondata, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "https://localhost:8080/home", bytes.NewReader(jsondata))
	cmd, filename, cleanup, err := GetCommand(req, true)
	if !asserts.Nil(err, err) {
		return
	}
	defer cleanup()
	f, err := os.ReadFile(filename)
	if !asserts.Nil(err, err) {
		return
	}

	if !asserts.Equal(string(f), string(jsondata), "body unmatch") {
		return
	}

	for i, val := range cmd {
		if val == "-d" {
			if !asserts.Equal(cmd[i+1], fmt.Sprintf("@%s", filename), "incorrect filename") {
				return
			}
			return
		}
	}
	asserts.Fail("unable get the body")
}

func TestGetResponse(t *testing.T) {
	asserts := assert.New(t)
	output := fmt.Sprintf(`%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n\r\n%s`,
		"HTTP/1.1 200 OK", "Connection: close", "X-Frame-Options: SAMEORIGIN",
		"Cache-Control: no-cache, no-store, must-revalidate",
		"Content-Length: 123", "Content-Type: application/json",
		"{\"user\":\"response\",\"flag\":false,\"count\":5}",
	)
	resp, err := GetResponse([]byte(output))
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

type testBody struct {
	User  string `json:"user,omitempty"`
	Flag  bool   `json:"flag,omitempty"`
	Count int    `json:"count,omitempty"`
}
