package z85

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
