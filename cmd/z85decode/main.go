package main

import (
	"bufio"
	"io"
	"os"

	"github.com/tortxof/z85"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	buf := make([]byte, 5)
	for {
		n, err := io.ReadFull(reader, buf[:5])
		if n == 5 {
			chunk := z85.Z85DecodeChunk([5]byte(buf))
			writer.Write(chunk[:])
		} else if n > 0 {
			writer.Write(z85.Z85Decode(buf[:n]))
		}
		if err != nil {
			break
		}
	}
}
