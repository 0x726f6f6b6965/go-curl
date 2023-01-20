# go-curl

[![GoDoc](https://godoc.org/github.com/0x726f6f6b6965/go-curl?status.svg)](https://godoc.org/github.com/0x726f6f6b6965/go-curl)
[![Go Report Card](https://goreportcard.com/badge/github.com/0x726f6f6b6965/go-curl)](https://goreportcard.com/report/github.com/0x726f6f6b6965/go-curl)
[![codecov](https://codecov.io/gh/0x726f6f6b6965/go-curl/branch/main/graph/badge.svg)](https://codecov.io/gh/0x726f6f6b6965/go-curl)

---

a repository for converting requests between curl and http

## Example

```go

httpReq, err := http.NewRequest(http.MethodGET, "https://example.com/", nil)

if err != nil {
    // Process error what you like
}

curlReq, err := NewCurlRequestWithContext(context.Background(), httpReq)

if err != nil {
    // Process error what you like
}


// if inscure is true, the command will add `-k`
err = curlReq.GenerateCommand(true)

if err != nil {
    // Process error what you like
}

// This can get the curl command
command := curlReq.GetCommands()

// This can run the curl command and get the http response

resp, err := curlReq.Do()

if err != nil {
    // Process error what you like
}


```
