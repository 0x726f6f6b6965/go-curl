package tocurl

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// GetCommand - return a string list which is the curl command
// the list can use strings.Join(result, " ") become a curl command string
func GetCommand(req *http.Request) ([]string, error) {
	cmd := []string{"curl", "-s", "-i"}

	if req.URL.Scheme == "https" {
		cmd = append(cmd, "-k")
	}

	for key, val := range req.Header {
		cmd = append(cmd, "-H", fmt.Sprintf("\"%s: %s\"", key, strings.Join(val, " ")))
	}

	cmd = append(cmd, "-X", req.Method)

	if req.Body != nil {
		buf, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("unable to read the request body, error: %w", err)
		}
		cmd = append(cmd, "-d", getCurlBody(buf))
	}

	cmd = append(cmd, req.URL.String())
	return cmd, nil
}

func GetResponse(buf []byte) (*http.Response, error) {
	result := &http.Response{}
	return result, nil
}

func getCurlBody(buf []byte) string {
	return fmt.Sprintf(`'%s'`, strings.Replace(string(buf), `'`, `'\''`, -1))
}
