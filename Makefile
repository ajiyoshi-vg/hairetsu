
all:

test:
	go test -v -cover ./...

verbose_test:
	go test -tags=verbose -v -cover ./...
