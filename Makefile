GO=/usr/bin/go
RM=/bin/rm
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
	${RM} -r bin/ pkg/
