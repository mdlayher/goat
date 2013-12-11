GO=/usr/bin/go
BIN=go-tracker
PATH=src/${BIN}/

${BIN}:
	${GO} build -o bin/${BIN} ${PATH}main.go

run:
	${GO} run ${PATH}main.go

clean:
	/bin/rm bin/${BIN}
