.PHONY: build clean deploy

build:
	env GOOS=linux go build  -ldflags="-s -w" -o bin/GetDevice devices/GetDevice/GetDevice.go
	env GOOS=linux go build  -ldflags="-s -w" -o bin/InsertDevice devices/InsertDevice/InsertDevice.go

clean:
	rm -rf bin 

deploy:  build
	sls deploy --verbose

format: 
	gofmt -w ./devices/devices.go