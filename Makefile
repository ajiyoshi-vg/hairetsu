
all:

test: generate
	go test -cover ./...

show_cover: generate
	go test -cover ./... -coverprofile=cover.out
	go tool cover -html=cover.out -o cover.html
	open cover.html

clean:
	-rm *.dat *.dict

generate:
	$(MAKE) -C doublearray $@

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

tiny_data:
	cat LICENSE| tr ' ' '\n' | grep -v "^$$" | sort -u | uniq > head.dat

head.dat: bench.dat Makefile
	tail -n 100000 $< > $@

%.trie: head.dat
	go run cmd/dump/main.go -o $@ -in $< -kind $*

bench: generate byte.trie rune.trie darts.trie dict.trie
	go test -benchmem -bench .

test_overhead: generate byte.trie
	go test -bench Overhead
