package tocurl

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
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
		cmd = append(cmd, "-H", fmt.Sprintf("%s: %s", key, strings.Join(val, " ")))
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
	result := &http.Response{
		Header: http.Header{},
	}
	sliceOutput := strings.Split(string(buf), "\r\n")
	if len(sliceOutput) < 1 {
		return nil, fmt.Errorf("unable to get the response")
	}

	if strings.Contains(sliceOutput[0], "HTTP") {
		code, status, err := getResponseStatus(sliceOutput[0])
		if err != nil {
			return nil, err
		}
		result.StatusCode = code
		result.Status = status
		sliceOutput = sliceOutput[1:]
	}

	var breakN int
	for i, val := range sliceOutput {
		if len(val) > 0 && strings.Contains(val, ":") {
			header := strings.SplitN(val, ":", 2)
			headerVals := strings.Split(header[1], ";")
			for _, headVal := range headerVals {
				result.Header.Add(header[0], headVal)
			}
		} else {
			breakN = i
			break
		}
	}

	if len(sliceOutput) > breakN+1 {
		body := sliceOutput[breakN+1]
		result.Body = io.NopCloser(strings.NewReader(body))
	}

	return result, nil
}

func getCurlBody(buf []byte) string {
	return fmt.Sprintf(`'%s'`, strings.Replace(string(buf), `'`, `'\''`, -1))
}

func getResponseStatus(output string) (int, string, error) {
	status := strings.Split(output, " ")
	code, err := strconv.Atoi(status[1])
	if err != nil {
		return 0, "", err
	}
	return code, http.StatusText(code), nil
}
