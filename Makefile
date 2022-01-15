
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

DONT_EDIT := // Code generated DO NOT EDIT

FILES := \
		 doublearray/mmap.words.go \
		 doublearray/mmap.bytes.go \
		 doublearray/mmap.runes.go \
		 doublearray/search.bytes.go \
		 doublearray/search.runes.go

generate: $(FILES)

SED := sed -i "" -e

clear_generated:
	rm $(FILES)

doublearray/mmap.words.go: doublearray/search.words.go
	echo "$(DONT_EDIT)" > $@
	cat $< >> $@
	$(SED) "s/Words/WordsMmap/" $@
	$(SED) "s/DoubleArray/Mmap/" $@
	goimports -w $@

%.bytes.go: %.words.go
	echo "$(DONT_EDIT)" > $@
	cat $< >> $@
	$(SED) "s/Words/Bytes/" $@
	$(SED) "s/cs word.Word/cs []byte/" $@
	$(SED) "s/Forward(c)/Forward(word.Code(c))/" $@
	goimports -w $@

%.runes.go: %.words.go
	echo "$(DONT_EDIT)" > $@
	cat $< >> $@
	$(SED) "s/struct{}/runedict.RuneDict/" $@
	$(SED) "s/(Words)/(s Runes)/" $@
	$(SED) "s/(WordsMmap)/(s RunesMmap)/" $@
	$(SED) "s/Words/Runes/" $@
	$(SED) "s/cs word.Word/cs string/" $@
	$(SED) "s/Forward(c)/Forward(s[c])/" $@
	goimports -w $@


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
	go test -benchmem -bench .

test_overhead: head.dat byte.dat
	go test -bench Overhead
