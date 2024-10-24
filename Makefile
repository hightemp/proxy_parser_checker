.PHONY: build clean

build:
	go build -o proxy_parser_checker ./cmd/main/main.go

clean:
	rm -f proxy_parser_checker

run: build
	./proxy_parser_checker	