package main

import (
	"strconv"
	"testing"
)

const n = '\n'

func BenchmarkBuffer(b *testing.B) {
	buf := make([]byte, 0, 2000000)

	for c := uint64(0); c < 1000000; c++ {
		buf = append(append(buf, strconv.FormatUint(c, 10)...), n)
	}
	_ = buf
}
