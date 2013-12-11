GO=/usr/bin/go
RM=/bin/rm
BIN=goat
DEP=bencode
PATH=src/
GOPATH=${PWD}

${BIN}:
	${GO} install ${BIN} ${DEP}
	${GO} build -o bin/${BIN} ${PATH}main.go

run:
	${GO} install ${BIN} ${DEP}
	${GO} run ${PATH}main.go

clean:
	${RM} -r bin/ pkg/
