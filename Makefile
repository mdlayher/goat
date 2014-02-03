make:
	go build -o bin/goat

ql:
	go build -tags='ql' -o bin/goat

fmt:
	go fmt
	go fmt github.com/mdlayher/goat/goat
	go fmt github.com/mdlayher/goat/goat/api
	go fmt github.com/mdlayher/goat/goat/common
	go fmt github.com/mdlayher/goat/goat/data
	go fmt github.com/mdlayher/goat/goat/data/udp
	go fmt github.com/mdlayher/goat/goat/tracker
	golint .
	golint goat
	golint goat/api
	golint goat/common
	golint goat/data
	golint goat/data/udp
	golint goat/tracker

test:
	go test github.com/mdlayher/goat/goat
	go test github.com/mdlayher/goat/goat/api
	go test github.com/mdlayher/goat/goat/common
	go test github.com/mdlayher/goat/goat/data
	go test -tags='ql' github.com/mdlayher/goat/goat/data
	go test github.com/mdlayher/goat/goat/data/udp
	go test github.com/mdlayher/goat/goat/tracker

darwin_386:
	GOOS="darwin" GOARCH="386" go build -o bin/goat_darwin_386

darwin_amd64:
	GOOS="darwin" GOARCH="amd64" go build -o bin/goat_darwin_amd64

linux_386:
	GOOS="linux" GOARCH="386" go build -o bin/goat_linux_386

linux_amd64:
	GOOS="linux" GOARCH="amd64" go build -o bin/goat_linux_amd64

windows_386:
	GOOS="windows" GOARCH="386" go build -o bin/goat_windows_386.exe

windows_amd64:
	GOOS="windows" GOARCH="amd64" go build -o bin/goat_windows_amd64.exe
