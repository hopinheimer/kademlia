package main

func XOR(a, b [20]byte) []byte {
	if len(a) != len(b) {
		panic("XOR: unequal lengths")
	}
	c := make([]byte, len(a))
	for i := range a {
		c[i] = a[i] ^ b[i]
	}
	return c
}

func PrefixLen(id []byte) int {
	for i := 0; i < IDLength; i++ {
		for j := 0; j < 8; j++ {
			if (id[i] & (0x80 >> j)) != 0 {
				return i*8 + j
			}
		}
	}
	return IDLength*8 - 1
}
