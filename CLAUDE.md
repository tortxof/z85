# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with
code in this repository.

## Interaction Style

Do not implement code or make changes directly, unless asked. Instead, answer
questions and provide explanations that help the user understand concepts well
enough to write implementations themselves. Act as a guide and teacher rather
than a code generator.

## Project Overview

This is a Go library implementing Z85 encoding/decoding (ZeroMQ Base-85). Z85
encodes binary data into printable ASCII characters using an 85-character
alphabet, with 4 bytes of input producing 5 bytes of output.

## Build Commands

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run a specific test
go test -run TestZ85Encode

# Run fuzz tests
go test -fuzz FuzzZ85RoundTrip
go test -fuzz FuzzStreamRoundTrip

# Run benchmarks
go test -bench=.

# Build CLI tools
go build ./cmd/z85encode
go build ./cmd/z85decode
```

## Architecture

- `z85.go` - Core library with `Z85Encode`/`Z85Decode` functions for byte slices,
  `Z85EncodeChunk`/`Z85DecodeChunk` for single 4-byte/5-byte chunks, and
  `NewEncoder`/`NewDecoder` for streaming `io.WriteCloser` implementations
- `z85_test.go` - Unit tests, fuzz tests, and benchmarks
- `cmd/z85encode/` - CLI tool that reads stdin and outputs Z85-encoded data
- `cmd/z85decode/` - CLI tool that reads Z85 stdin and outputs decoded binary

## Encoding Details

- Input is processed in 4-byte chunks (encode) or 5-byte chunks (decode)
- Partial chunks are padded internally; output is trimmed to exclude padding
  bytes
- The `#` character is used for padding during decode
- Lookup tables (`decodeMap`, `decodeMultipliers`) are initialized at package
  load time for performance
