
all:

test:
	go test -cover ./...

show_cover:
	go test -cover ./... -coverprofile=cover.out
	go tool cover -html=cover.out -o cover.html
	open cover.html

clean:
	-rm *.dat *.dict

URL := https://dumps.wikimedia.org/jawiki/latest/jawiki-latest-all-titles.gz

jawiki-latest-all-titles.gz:
	# donwloading wikipedia article titles.
	# see their license
	#   https://dumps.wikimedia.org/legal.html
	# CC BY-SA 3.0
	#   https://creativecommons.org/licenses/by-sa/3.0/
	curl -o $@ $(URL)

bench.dat: jawiki-latest-all-titles.gz
	gunzip -c $< | cut -f 2 | sort | uniq > $@

head.dat: Makefile
	tail -n 100000 bench.dat > $@

byte.dat: head.dat
	go run cmd/dump/main.go -o $@ -in $< -kind byte

rune.dat rune.dat.dict: head.dat
	go run cmd/dump/main.go -o $@ -in $< -kind rune

darts.dat : head.dat
	go run cmd/dump/main.go -o $@ -in $< -kind darts

bench: head.dat byte.dat rune.dat rune.dat.dict darts.dat
	go test -bench .
