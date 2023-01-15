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
)

func TestGetCommandWithBasicAuth(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://localhost:8080/home", nil)
	user, pwd := "admin", "admin123"
	req.SetBasicAuth(user, pwd)
	req.BasicAuth()
	cmd, _, _, err := GetCommand(req, true)
	if err != nil {
		t.Error(err)
		return
	}
	for i, val := range cmd {
		if strings.Contains(val, "Authorization") {
			if cmd[i-1] != "-H" {
				t.Error("incorrect header format")
				return
			}
			auths := strings.Split(val, ":")
			if len(auths) != 2 {
				t.Error("incorrect header value")
				return
			}
			if !strings.HasPrefix(auths[1], " Basic ") {
				t.Error("incorrect prefix of the header value")
				return
			}
			c, err := base64.StdEncoding.DecodeString(strings.Split(auths[1], " ")[2])
			if err != nil {
				t.Error("unable decode the header value")
				return
			}
			cs := string(c)
			username, password, ok := strings.Cut(cs, ":")
			if !ok {
				t.Error("unable parse the auth")
				return
			}
			if username != user || password != pwd {
				t.Error("incorrect username/password")
				return
			}
			return
		}
	}
	t.Error("unable get the auth")
}

func TestGetCommandWithBody(t *testing.T) {
	body := &testBody{
		User:  "admin",
		Flag:  true,
		Count: 3,
	}
	jsondata, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "https://localhost:8080/home", bytes.NewReader(jsondata))
	cmd, filename, cleanup, err := GetCommand(req, true)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()
	f, err := os.ReadFile(filename)
	if err != nil {
		t.Error(err)
		return
	}

	if string(f) != string(jsondata) {
		t.Error(err)
		return
	}

	for i, val := range cmd {
		if val == "-d" {
			if cmd[i+1] != fmt.Sprintf("@%s", filename) {
				t.Error("incorrect filename")
				return
			}
			return
		}
	}
	t.Error("unable get the body")
	return
}

func TestGetResponse(t *testing.T) {
	output := fmt.Sprintf(`%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n%s\r\n\r\n%s`,
		"HTTP/1.1 200 OK", "Connection: close", "X-Frame-Options: SAMEORIGIN",
		"Cache-Control: no-cache, no-store, must-revalidate",
		"Content-Length: 123", "Content-Type: application/json",
		"{\"user\":\"response\",\"flag\":false,\"count\":5}",
	)
	resp, err := GetResponse([]byte(output))
	if err != nil {
		t.Error(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		t.Error("incorrect status code")
		return
	}

	if resp.Status != http.StatusText(http.StatusOK) {
		t.Error("incorrect status code")
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error("read body error")
		return
	}
	var result testBody
	err = json.Unmarshal(body, &result)
	if err != nil {
		t.Error("parse body error")
		return
	}
	if result.Count != 5 || result.Flag != false || result.User != "response" {
		t.Error("incorrect body")
		return
	}
}

type testBody struct {
	User  string `json:"user,omitempty"`
	Flag  bool   `json:"flag,omitempty"`
	Count int    `json:"count,omitempty"`
}
