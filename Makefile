BINARY_NAME=aws_restore_glacier

build:
	mkdir -p bin/linux
	mkdir -p bin/windows
	mkdir -p bin/darwin
	GOOS=linux GOARCH=amd64 go build -o bin/linux/$(BINARY_NAME)
	GOOS=windows GOARCH=amd64 go build -o bin/windows/$(BINARY_NAME).exe
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin/$(BINARY_NAME)
