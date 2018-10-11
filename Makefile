SHELL := /bin/bash
.SILENT:
help:
	echo
	echo "Mobius Golang Client Make Commands"
	echo
	echo "  Commands  "
	echo
	echo "    help - show this message  "
	echo "    example-auth -  build/run the auth example"
	echo "    example-app  -  build/run the flappy app example"
	echo "    test - run the tests for this client"
	echo "    test-verbose - run the tests, and show coverage, and race conditions"
	echo "    deps - check status of dependencies"


test-verbose:
	ginkgo -v -cover ./...

test:
	ginkgo -r --cover --progress --succinct

deps:
	dep ensure

example-auth:
		source ./.env.example
		cd examples/auth && go run main.go

example-app:
		source ./.env.example
		cd examples/flappy && go run main.go
