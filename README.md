# ⚡ z85

A Go library for Z85 encoding and decoding, implementing the [ZeroMQ Base-85
encoding](https://rfc.zeromq.org/spec/32/).

Z85 encodes binary data into printable ASCII characters using an 85-character
alphabet. Every 4 bytes of input produce 5 bytes of output, resulting in a 25%
size increase (compared to 33% for Base64).

## Installation

```bash
go get github.com/tortxof/z85
```

## Library Usage

```go
package main

import (
	"fmt"

	"github.com/tortxof/z85"
)

func main() {
	// Encode binary data
	data := []byte{0x86, 0x4F, 0xD2, 0x6F, 0xB5, 0x59, 0xF7, 0x5B}
	encoded := z85.Z85Encode(data)
	fmt.Println(string(encoded)) // HelloWorld

	// Decode Z85 string
	decoded := z85.Z85Decode([]byte("HelloWorld"))
	fmt.Printf("%x\n", decoded) // 864fd26fb559f75b
}
```

## API

### Byte Slice Functions

- `Z85Encode(data []byte) []byte` - Encode a byte slice to Z85
- `Z85Decode(data []byte) []byte` - Decode a Z85 byte slice to binary
- `Z85EncodeChunk(chunk [4]byte) [5]byte` - Encode a single 4-byte chunk
- `Z85DecodeChunk(chunk [5]byte) [4]byte` - Decode a single 5-byte chunk

Partial chunks are handled automatically with internal padding.

### Streaming API

- `NewEncoder(w io.Writer) io.WriteCloser` - Create a streaming encoder
- `NewDecoder(w io.Writer) io.WriteCloser` - Create a streaming decoder

The streaming types implement `io.WriteCloser`. Data written to the
encoder/decoder is processed and written to the underlying writer. Call
`Close()` to flush any buffered partial chunks.

```go
// Streaming encode
var buf bytes.Buffer
encoder := z85.NewEncoder(&buf)
encoder.Write([]byte{0x86, 0x4F, 0xD2, 0x6F})
encoder.Write([]byte{0xB5, 0x59, 0xF7, 0x5B})
encoder.Close()
fmt.Println(buf.String()) // HelloWorld

// Streaming decode
var out bytes.Buffer
decoder := z85.NewDecoder(&out)
io.Copy(decoder, strings.NewReader("HelloWorld"))
decoder.Close()
fmt.Printf("%x\n", out.Bytes()) // 864fd26fb559f75b
```

## CLI Tools

Build the command-line tools:

```bash
go build ./cmd/z85encode
go build ./cmd/z85decode
```

Usage:

```bash
# Encode
echo -n "test" | ./z85encode
# Output: By/Jn

# Decode
echo -n "By/Jn" | ./z85decode
# Output: test

# Encode a file
./z85encode < input.bin > output.z85

# Decode a file
./z85decode < input.z85 > output.bin
```

## Testing

```bash
# Run tests
go test

# Run benchmarks
go test -bench=.

# Run fuzz tests
go test -fuzz=FuzzZ85RoundTrip
go test -fuzz=FuzzStreamRoundTrip
```

## License

See [LICENSE](LICENSE) for details.
