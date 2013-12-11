GO=/usr/bin/go
BIN=goat
PATH=src/

${BIN}:
	${GO} build -o bin/${BIN} ${PATH}main.go

run:
	${GO} run ${PATH}main.go

clean:
	/bin/rm bin/${BIN}
