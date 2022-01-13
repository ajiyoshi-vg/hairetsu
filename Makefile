
all:

test:
	go test -v -cover ./...

show_cover: prepare
	open cover.html

prepare:
	go test -cover ./... -coverprofile=cover.out
	go tool cover -html=cover.out -o $@
