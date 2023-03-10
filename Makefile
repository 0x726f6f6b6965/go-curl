PROJECTNAME := $(shell basename "$(PWD)")
GOPATH := $(shell go env GOPATH)

.PHONY: default help test clean-test
default: help

## test: run all test case in the project
test:
	@ go test -coverprofile cover.out ./...
	@ go tool cover -html=cover.out -o cover.html 
	@ open cover.html

## clean-test: clean up the files which generated by test commands
clean-test:
	@ rm cover.html
	@ rm cover.out

## generated-file-type: use stringer to generate the constant value of file type
generated-file-type:
	@ $(GOPATH)/bin/stringer -type=FileType ./filetype/file_type.go

## install-stringer: install stringer
install-stringer:
	@ go get -u golang.org/x/tools/cmd/stringer

# help
help: Makefile
	@echo
	@echo "Choose a command to run in $(PROJECTNAME):"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo
