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

	buf := make([]byte, 4)
	for {
		n, err := io.ReadFull(reader, buf[:4])
		if n == 4 {
			chunk := z85.Z85EncodeChunk([4]byte(buf))
			writer.Write(chunk[:])
		} else if n > 0 {
			writer.Write(z85.Z85Encode(buf[:n]))
		}
		if err != nil {
			break
		}
	}
}
