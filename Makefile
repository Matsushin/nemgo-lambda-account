.PHONY: deps clean build

deps:
	go get -u ./...

clean: 
	rm -rf ./account/account
	
build:
	GOOS=linux GOARCH=amd64 go build -o account/account ./account