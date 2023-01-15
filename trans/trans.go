package trans

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// GetCommand - return a string list which is the curl command
// and a function which can remove the temporary file of request body.
// if inscure is true, the command will add '-k'
// the list can use strings.Join(result, " ") become a curl command string
func GetCommand(req *http.Request, inscure bool) ([]string, string, func() error, error) {
	cmd := []string{"curl", "-s", "-i"}
	cleanup := func() error { return nil }
	if inscure {
		cmd = append(cmd, "-k")
	}

	for key, val := range req.Header {
		cmd = append(cmd, "-H", fmt.Sprintf("%s: %s", key, strings.Join(val, " ")))
	}

	cmd = append(cmd, "-X", req.Method)

	var filename string
	if req.Body != nil {

		buf, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, "", cleanup,
				fmt.Errorf("unable to read the request body, error: %w", err)
		}
		// create a temporary file to avoid
		// the error of argument list too long
		filename = uuid.New().String()
		err = ioutil.WriteFile(filename, buf, 0644)
		if err != nil {
			return nil, "", cleanup,
				fmt.Errorf("unable to create the file of request body, error: %w", err)
		}
		cmd = append(cmd, "-d", fmt.Sprintf("@%s", filename))
	}

	cmd = append(cmd, req.URL.String())
	cleanup = func() error {
		if len(filename) == 0 {
			return nil
		}
		return os.Remove(filename)
	}
	return cmd, filename, cleanup, nil
}

// Execute - execute the curl command
func Execute(command []string) ([]byte, error) {
	if len(command) < 3 && command[0] != "curl" {
		return nil, fmt.Errorf("unable to run this command, command: %s", command)
	}
	return exec.Command(command[0], command[1:]...).CombinedOutput()
}

// Execute - execute the curl command with context
func ExecuteWithContext(ctx context.Context, command []string) ([]byte, error) {
	if len(command) < 3 && command[0] != "curl" {
		return nil, fmt.Errorf("unable to run this command, command: %s", command)
	}
	return exec.CommandContext(ctx, command[0], command[1:]...).CombinedOutput()
}

// GetResponse - translate the command to a http response
func GetResponse(buf []byte) (*http.Response, error) {
	result := &http.Response{
		Header: http.Header{},
	}
	sliceOutput := strings.Split(string(buf), "\\r\\n")
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
			result.Header.Add(header[0], strings.TrimLeft(header[1], " "))
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

func getResponseStatus(output string) (int, string, error) {
	status := strings.Split(output, " ")
	code, err := strconv.Atoi(status[1])
	if err != nil {
		return 0, "", err
	}
	return code, http.StatusText(code), nil
}
