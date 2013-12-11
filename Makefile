GO=/usr/bin/go
BIN=goat
PATH=src/
GOPATH=${PWD}

${BIN}:
	${GO} install ${BIN}
	${GO} build -o bin/${BIN} ${PATH}main.go

run:
	${GO} install ${BIN}
	${GO} run ${PATH}main.go

clean:
	/bin/rm bin/${BIN}
