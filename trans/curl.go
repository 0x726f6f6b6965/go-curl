package trans

import (
	"context"
	"net/http"

	"github.com/0x726f6f6b6965/go-curl/filetype"
)

type curlRequest struct {
	ctx            context.Context
	request        *http.Request
	cmd            []string
	filename       string
	cleanup        func() error
	private        *privatekey
	cert           *certificate
	ca             string
	combinedOutput func(ctx context.Context, name string, args ...string) ([]byte, error)
}

type CurlRequest interface {
	// GenerateCommand - generate curl command with curl request
	GenerateCommand(inscure bool) error

	// GetCommands - get the curl command
	GetCommands() []string

	// GetFilename - get the filename which contains the request body
	GetFilename() string

	// Execute - execute the curl command
	Execute() ([]byte, error)

	// Do - execute the curl command and get the http.Response, like http.client.Do
	Do() (*http.Response, error)

	// AddHeader - add curl header
	AddHeader(key, value string) error

	// GetHeaders - get all the curl headers
	GetHeaders() http.Header

	// GetHeader - get the specific header value
	GetHeader(key string) []string

	// SetPrivateKey - set the private key
	SetPrivateKey(fileType filetype.FileType, path string)

	// SetCertificate - set the certificate
	SetCertificate(fileType filetype.FileType, path, password string)

	// SetCA - set the ca file
	SetCA(filePath string)
}

type privatekey struct {
	fileType string
	path     string
}

type certificate struct {
	fileType string
	path     string
	password string
}
