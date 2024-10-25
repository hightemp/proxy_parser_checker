.PHONY: build clean

build:
	go build -o proxy_parser_checker ./cmd/main/main.go

build_static:
	CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o proxy_parser_checker_static ./cmd/main/main.go

clean:
	rm -f proxy_parser_checker

run: build
	./proxy_parser_checker	