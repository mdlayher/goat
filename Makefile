BIN=go-tracker

${BIN}:
	go build -o ${BIN} src/main.go

run:
	go run src/main.go

clean:
	rm ${BIN}
