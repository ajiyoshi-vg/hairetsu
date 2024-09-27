package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/google/uuid"
)

var opt struct {
	num int
}

func init() {
	flag.IntVar(&opt.num, "n", 1000*1000, "number")
	flag.Parse()
}
func main() {
	if err := run(); err != nil {
		slog.Error(err.Error())
	}
}

func run() error {
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	for range opt.num {
		fmt.Fprintln(w, uuid.New().String())
	}

	return nil
}
