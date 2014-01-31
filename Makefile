make:
	go build -o bin/goat

fmt:
	go fmt
	go fmt github.com/mdlayher/goat/goat
	golint .
	golint goat

test:
	go test github.com/mdlayher/goat/goat

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
