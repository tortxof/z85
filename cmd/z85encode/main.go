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
	encoder := z85.NewEncoder(writer)
	defer encoder.Close()
	io.Copy(encoder, os.Stdin)
}
