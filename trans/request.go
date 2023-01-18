package trans

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/0x726f6f6b6965/go-curl/filetype"
)

func NewCurlRequestWithContext(ctx context.Context, req *http.Request) (CurlRequest, error) {
	if req == nil {
		return nil, fmt.Errorf("unable create a curl request without http.Request")
	}

	return &curlRequest{ctx: ctx, request: req,
		cleanup: func() error { return nil },
		combinedOutput: func(ctx context.Context, name string, args ...string) ([]byte, error) {
			return exec.CommandContext(ctx, name, args...).CombinedOutput()
		}}, nil
}

func NewCurlRequest(req *http.Request) (CurlRequest, error) {
	return NewCurlRequestWithContext(context.Background(), req)
}

func (curl *curlRequest) GenerateCommand(inscure bool) error {
	cmd, filename, cleanup, err := getCommand(curl.request, inscure, curl.private, curl.cert, curl.ca)
	if err != nil {
		return err
	}

	curl.cmd = cmd
	curl.filename = filename
	curl.cleanup = cleanup
	return nil
}

func (curl *curlRequest) GetCommands() []string {
	return curl.cmd
}

func (curl *curlRequest) GetFilename() string {
	return curl.filename
}

func (curl *curlRequest) Execute() ([]byte, error) {
	defer curl.cleanup()
	if len(curl.cmd) < 3 {
		return nil, fmt.Errorf("unable to run this command, command: %s", curl.cmd)
	}
	return curl.combinedOutput(curl.ctx, curl.cmd[0], curl.cmd[1:]...)
}

func (curl *curlRequest) Do() (*http.Response, error) {
	resp, err := curl.Execute()
	if err != nil {
		return nil, err
	}
	return getResponse(resp)
}

func (curl *curlRequest) AddHeader(key, value string) error {
	curl.request.Header.Add(key, value)

	// check the GenerateCommand was run.
	if curl.cmd != nil {
		if curl.filename != "" {
			err := curl.cleanup()
			if err != nil {
				return err
			}
		}
		if curl.cmd[3] != "-k" {
			return curl.GenerateCommand(false)
		}
		return curl.GenerateCommand(true)
	}

	return nil
}

func (curl *curlRequest) GetHeaders() http.Header {
	header := curl.request.Header
	return header
}

func (curl *curlRequest) GetHeader(key string) []string {
	values := curl.request.Header.Values(key)
	return values
}

func (curl *curlRequest) SetPrivateKey(fileType filetype.FileType, path string) {
	curl.private = &privatekey{fileType: fileType.String(), path: path}
}

func (curl *curlRequest) SetCertificate(fileType filetype.FileType, path, password string) {
	curl.cert = &certificate{fileType: fileType.String(), path: path, password: password}
}

func (curl *curlRequest) SetCA(filePath string) {
	curl.ca = filePath
}
