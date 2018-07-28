package main

import (
	"math"
	"os"
	"runtime"
	"strconv"
)

const (
	_  = iota
	KB = 1 << (10 * iota)
	MB = 1 << (10 * iota)
)

const (
	chanSize = KB
	bufSize  = MB
)

func main() {
	thread := uint64(runtime.NumCPU())
	c := make(chan uint64, chanSize)
	for i := uint64(3); i < uint64(2*thread+3); i += 2 {
		go prime(i, thread, c)
	}
	buf := make([]byte, 0)
	for {
		buf = append(buf, strconv.FormatUint(<-c, 10)...)
		buf = append(buf, '\n')
		if len(buf) > bufSize {
			os.Stdout.Write(buf)
			buf = make([]byte, 0)
		}
	}
}

func prime(n uint64, thread uint64, c chan uint64) {
	increment := uint64(2 * thread)
	max := uint64(math.MaxUint64 - increment)
	for i := n; i < max; i += increment {
		if isPrime(i) {
			c <- i
		}
	}
}

func isPrime(i uint64) bool {
	for j := uint64(3); j*j <= i; j += uint64(2) {
		if i%j == 0 {
			return false
		}
	}
	return true
}
