package main

import (
	"log"
	"math"
	"math/big"
	"os"
	"runtime"
	"strconv"
	"sync"
)

const (
	limit     = uint64(math.MaxUint32)
	calcRange = 10000000
)

type SafeNum struct {
	n uint64
	m sync.RWMutex
}

func (s *SafeNum) Get() uint64 {
	s.m.Lock()
	tmp := s.n

	s.n += calcRange

	s.m.Unlock()
	return tmp
}

func main() {
	s := &SafeNum{ /*n: 2*/ }
	thread := uint64(runtime.NumCPU())

	wg := &sync.WaitGroup{}

	for i := uint64(0); i < thread; i++ {
		wg.Add(1)
		go prime(s, wg)
	}
	wg.Wait()
}

const (
	bufferSize  = 20000000
	bufferLimit = 19999000
)

func prime(s *SafeNum, wg *sync.WaitGroup) {

	buf := make([]byte, 0, bufferSize)
	m := &sync.Mutex{}
	for num := s.Get(); num <= limit; num = s.Get() {

		for i := uint64(1); i != calcRange; i++ {
			num += 1
			// 多段チェックにすることで高速化している
			if isPrime(num) && big.NewInt(int64(num)).ProbablyPrime(0) {
				buf = append(append(buf, strconv.FormatUint(num, 10)...), '\n')
				if len(buf) > bufferLimit {

					_, err := os.Stdout.Write(buf)
					if err != nil {
						log.Fatal(err)
					}
					buf = make([]byte, 0, bufferSize)

				}
				m.Unlock()
			}
		}
	}
	_, err := os.Stdout.Write(buf)
	if err != nil {
		log.Fatal(err)
	}
	wg.Done()
}

func isPrime(n uint64) bool {

	// bases of 2, 7, 61 are sufficient to cover 2^32
	switch n {
	case 0, 1:
		return false
	case 2, 7, 61:
		return true
	}
	// compute s, d where 2^s * d = n-1
	nm1 := n - 1
	d := nm1
	s := 0
	for d&1 == 0 {
		d >>= 1
		s++
	}
	n64 := uint64(n)
	for _, a := range []uint32{2, 7, 61} {
		// compute x := a^d % n
		x := uint64(1)
		p := uint64(a)
		for dr := d; dr > 0; dr >>= 1 {
			if dr&1 != 0 {
				x = x * p % n64
			}
			p = p * p % n64
		}
		if x == 1 || uint64(x) == nm1 {
			continue
		}
		for r := 1; ; r++ {
			if r >= s {
				return false
			}
			x = x * x % n64
			if x == 1 {
				return false
			}
			if uint64(x) == nm1 {
				break
			}
		}
	}
	return true
}
