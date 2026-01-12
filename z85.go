package z85

import "io"

const Z85Alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.-:+=^!/*?&<>()[]{}@%$#"

var decodeMap [256]int

var decodeMultipliers [5]uint32

var paddingChunk = [5]byte{'#', '#', '#', '#', '#'}

func init() {
	for i := range Z85Alphabet {
		decodeMap[Z85Alphabet[i]] = i
	}
	var multiplier uint32 = 85 * 85 * 85 * 85
	for i := range 5 {
		decodeMultipliers[i] = multiplier
		multiplier /= 85
	}
}

// Process one 4 byte chunk of input data. Return 5 bytes of encoded data.
func Z85EncodeChunk(chunk [4]byte) [5]byte {
	value := uint32(chunk[0])<<24 | uint32(chunk[1])<<16 | uint32(chunk[2])<<8 | uint32(chunk[3])

	result := [5]byte{}

	for i := 4; i >= 0; i-- {
		result[i] = Z85Alphabet[value%85]
		value /= 85
	}

	return result
}

func Z85Encode(data []byte) []byte {
	numInputChunks := len(data) / 4
	remainingBytes := len(data) % 4
	padding := 0
	if remainingBytes > 0 {
		numInputChunks++
		padding = 4 - remainingBytes
	}
	result := make([]byte, numInputChunks*5)
	for chunkNum := range numInputChunks {
		var chunk [4]byte
		copy(chunk[:], data[chunkNum*4:])
		encoded := Z85EncodeChunk(chunk)
		copy(result[chunkNum*5:], encoded[:])
	}
	return result[:len(result)-padding]
}

// Process one 5 byte chunk of Z85 encoded data. Return 4 bytes of decoded data.
func Z85DecodeChunk(chunk [5]byte) [4]byte {
	var value uint32
	var result = [4]byte{}

	for i, multiplier := range decodeMultipliers {
		value += uint32(decodeMap[chunk[i]]) * multiplier
	}

	for i := range 4 {
		result[i] = byte(value >> (24 - (i * 8)))
	}

	return result
}

func Z85Decode(data []byte) []byte {
	numInputChunks := len(data) / 5
	remainingBytes := len(data) % 5
	padding := 0
	if remainingBytes > 0 {
		numInputChunks++
		padding = 5 - remainingBytes
	}
	result := make([]byte, numInputChunks*4)
	for chunkNum := range numInputChunks {
		chunk := paddingChunk
		copy(chunk[:], data[chunkNum*5:])
		decoded := Z85DecodeChunk(chunk)
		copy(result[chunkNum*4:], decoded[:])
	}
	return result[:len(result)-padding]
}

type Encoder struct {
	w   io.Writer
	buf [4]byte
	n   int
	err error
}

func (enc *Encoder) writeChunk(chunk [4]byte) {
	encoded := Z85EncodeChunk(chunk)
	_, err := enc.w.Write(encoded[:])
	if err != nil {
		enc.err = err
	}
}

func (enc *Encoder) Close() error {
	if enc.n > 0 {
		padding := 4 - enc.n
		chunk := [4]byte{}
		copy(chunk[:], enc.buf[:enc.n])
		encodedChunk := Z85EncodeChunk(chunk)
		_, err := enc.w.Write(encodedChunk[:5-padding])
		if err != nil {
			enc.err = err
		}
	}
	return enc.err
}

func (enc *Encoder) Write(p []byte) (int, error) {
	if enc.err != nil {
		return 0, enc.err
	}

	bytesIn := len(p)

	if enc.n > 0 {
		bytesNeeded := 4 - enc.n
		if bytesIn >= bytesNeeded {
			copy(enc.buf[enc.n:], p[:bytesNeeded])
			enc.writeChunk(enc.buf)
			p = p[bytesNeeded:]
			enc.n = 0
		} else {
			copy(enc.buf[enc.n:], p)
			enc.n += bytesIn
			return bytesIn, nil
		}
	}

	if enc.err != nil {
		return 0, enc.err
	}

	for len(p) >= 4 {
		var chunk [4]byte
		copy(chunk[:], p[:4])
		enc.writeChunk(chunk)
		if enc.err != nil {
			return 0, enc.err
		}
		p = p[4:]
	}

	enc.n = copy(enc.buf[:], p)

	return bytesIn, enc.err
}

func NewEncoder(w io.Writer) io.WriteCloser {
	return &Encoder{w: w}
}

type Decoder struct {
	w   io.Writer
	buf [5]byte
	n   int
	err error
}

func (dec *Decoder) writeChunk(chunk [5]byte) {
	decoded := Z85DecodeChunk(chunk)
	_, err := dec.w.Write(decoded[:])
	if err != nil {
		dec.err = err
	}
}

func (dec *Decoder) Close() error {
	if dec.n > 0 {
		padding := 5 - dec.n
		chunk := paddingChunk
		copy(chunk[:], dec.buf[:dec.n])
		decodedChunk := Z85DecodeChunk(chunk)
		_, err := dec.w.Write(decodedChunk[:4-padding])
		if err != nil {
			dec.err = err
		}
	}
	return dec.err
}

func (dec *Decoder) Write(p []byte) (int, error) {
	if dec.err != nil {
		return 0, dec.err
	}

	bytesIn := len(p)

	if dec.n > 0 {
		bytesNeeded := 5 - dec.n
		if bytesIn >= bytesNeeded {
			copy(dec.buf[dec.n:], p[:bytesNeeded])
			dec.writeChunk(dec.buf)
			p = p[bytesNeeded:]
			dec.n = 0
		} else {
			copy(dec.buf[dec.n:], p)
			dec.n += bytesIn
			return bytesIn, nil
		}
	}

	if dec.err != nil {
		return 0, dec.err
	}

	for len(p) >= 5 {
		var chunk [5]byte
		copy(chunk[:], p[:5])
		dec.writeChunk(chunk)
		if dec.err != nil {
			return 0, dec.err
		}
		p = p[5:]
	}

	dec.n = copy(dec.buf[:], p)

	return bytesIn, dec.err
}

func NewDecoder(w io.Writer) io.WriteCloser {
	return &Decoder{w: w}
}
