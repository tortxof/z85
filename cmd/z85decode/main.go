package main

import (
	"bufio"
	"io"
	"os"

	"github.com/tortxof/z85"
)

func main() {
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	decoder := z85.NewDecoder(writer)
	defer decoder.Close()
	io.Copy(decoder, os.Stdin)
}
