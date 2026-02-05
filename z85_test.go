package z85

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"testing"
)

func TestEncodeChunk(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input [4]byte
		want  [5]byte
	}{
		{
			name:  "zero bytes",
			input: [4]byte{0x00, 0x00, 0x00, 0x00},
			want:  [5]byte{'0', '0', '0', '0', '0'},
		},
		{
			name:  "known value",
			input: [4]byte{0x86, 0x4F, 0xD2, 0x6F},
			want:  [5]byte{'H', 'e', 'l', 'l', 'o'},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeChunk(tt.input)
			if got != tt.want {
				t.Errorf("EncodeChunk(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

var encodeTests = []struct {
	name  string
	input []byte
	want  string
}{
	{
		name:  "8 bytes",
		input: []byte{0x86, 0x4F, 0xD2, 0x6F, 0xB5, 0x59, 0xF7, 0x5B},
		want:  "HelloWorld",
	},
	{
		name:  "text with padding",
		input: []byte("Hello world!!"),
		want:  "nm=QNzY<mxA+]nfaP",
	},
	{
		name:  "empty input",
		input: []byte{},
		want:  "",
	},
	{
		name:  "single byte",
		input: []byte{0x42},
		want:  "li",
	},
}

func TestEncode(t *testing.T) {
	t.Parallel()

	for _, tt := range encodeTests {
		t.Run(tt.name, func(t *testing.T) {
			got := Encode(tt.input)
			if string(got) != tt.want {
				t.Errorf("Encode(%v) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

var roundTripTests = []struct {
	name  string
	input []byte
}{
	{
		name:  "8 bytes",
		input: []byte{0x86, 0x4F, 0xD2, 0x6F, 0xB5, 0x59, 0xF7, 0x5B},
	},
	{
		name:  "text",
		input: []byte("Hello world!!"),
	},
	{
		name:  "empty",
		input: []byte{},
	},
	{
		name:  "single byte",
		input: []byte{0x42},
	},
	{
		name:  "all zeros",
		input: []byte{0x00, 0x00, 0x00, 0x00},
	},
	{
		name:  "all ones",
		input: []byte{0xFF, 0xFF, 0xFF, 0xFF},
	},
}

func TestZ85RoundTrip(t *testing.T) {
	t.Parallel()

	for _, tt := range roundTripTests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := Encode(tt.input)
			decoded := Decode(encoded)
			if !bytes.Equal(decoded, tt.input) {
				t.Errorf("Round trip failed: input %v, encoded %q, decoded %v", tt.input, encoded, decoded)
			}
		})
	}
}

func FuzzZ85RoundTrip(f *testing.F) {
	f.Add([]byte{0x86, 0x4F, 0xD2, 0x6F, 0xB5, 0x59, 0xF7, 0x5B})
	f.Add([]byte("Hello world!!"))
	f.Add([]byte{})
	f.Add([]byte{0x42})

	f.Fuzz(func(t *testing.T, input []byte) {
		encoded := Encode(input)
		decoded := Decode(encoded)

		if !bytes.Equal(input, decoded) {
			t.Errorf("Round trip failed: input %v (len=%d), encoded %q (len=%d), decoded %v (len=%d)",
				input, len(input), encoded, len(encoded), decoded, len(decoded))
		}
	})
}

func BenchmarkEncode(b *testing.B) {
	sizes := []int{16, 256, 4096, 16384}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			input := make([]byte, size)
			rand.Read(input)
			b.ResetTimer()
			b.SetBytes(int64(size))

			for i := 0; i < b.N; i++ {
				Encode(input)
			}
		})
	}
}

func BenchmarkDecode(b *testing.B) {
	sizes := []int{16, 256, 4096, 16384}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			input := make([]byte, size)
			rand.Read(input)
			encoded := Encode(input)
			b.ResetTimer()
			b.SetBytes(int64(size))

			for i := 0; i < b.N; i++ {
				Decode(encoded)
			}
		})
	}
}

func TestStreamEncode(t *testing.T) {
	t.Parallel()

	for _, tt := range encodeTests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			encoder := NewEncoder(&buf)
			encoder.Write(tt.input)
			err := encoder.Close()
			if err != nil {
				t.Fatalf("Close() returned an unexpected error: %v", err)
			}
			got := buf.String()
			if string(got) != tt.want {
				t.Errorf("Encoder(%v) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestStreamRoundTrip(t *testing.T) {
	t.Parallel()

	for _, tt := range roundTripTests {
		t.Run(tt.name, func(t *testing.T) {
			var encBuf bytes.Buffer
			encoder := NewEncoder(&encBuf)
			encoder.Write(tt.input)
			err := encoder.Close()
			if err != nil {
				t.Fatalf("Close() returned an unexpected error: %v", err)
			}
			encoded := encBuf.Bytes()
			var decBuf bytes.Buffer
			decoder := NewDecoder(&decBuf)
			decoder.Write(encoded)
			err = decoder.Close()
			if err != nil {
				t.Fatalf("Close() returned an unexpected error: %v", err)
			}
			decoded := decBuf.Bytes()
			if !bytes.Equal(decoded, tt.input) {
				t.Errorf("Round trip failed: input %v, encoded %q, decoded %v", tt.input, encoded, decoded)
			}
		})
	}
}

func FuzzStreamRoundTrip(f *testing.F) {
	f.Add([]byte{0x86, 0x4F, 0xD2, 0x6F, 0xB5, 0x59, 0xF7, 0x5B}, uint8(1))
	f.Add([]byte("Hello world!!"), uint8(3))
	f.Add([]byte{}, uint8(1))
	f.Add([]byte{0x42}, uint8(7))
	f.Add([]byte{0xFF, 0xFF, 0xFF, 0xFF}, uint8(4))

	f.Fuzz(func(t *testing.T, input []byte, chunkSize uint8) {
		if chunkSize == 0 {
			chunkSize = 1
		}

		// Encode in chunks
		var encBuf bytes.Buffer
		encoder := NewEncoder(&encBuf)
		data := input
		for len(data) > 0 {
			n := int(chunkSize)
			if n > len(data) {
				n = len(data)
			}
			encoder.Write(data[:n])
			data = data[n:]
		}
		if err := encoder.Close(); err != nil {
			t.Fatalf("encoder.Close() error: %v", err)
		}
		encoded := encBuf.Bytes()

		// Decode in chunks
		var decBuf bytes.Buffer
		decoder := NewDecoder(&decBuf)
		data = encoded
		for len(data) > 0 {
			n := int(chunkSize)
			if n > len(data) {
				n = len(data)
			}
			decoder.Write(data[:n])
			data = data[n:]
		}
		if err := decoder.Close(); err != nil {
			t.Fatalf("decoder.Close() error: %v", err)
		}
		decoded := decBuf.Bytes()

		if !bytes.Equal(input, decoded) {
			t.Errorf("Round trip failed: input %v (len=%d), encoded %q (len=%d), decoded %v (len=%d), chunkSize=%d",
				input, len(input), encoded, len(encoded), decoded, len(decoded), chunkSize)
		}
	})
}
