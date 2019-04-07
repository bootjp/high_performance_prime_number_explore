package main

import (
	"log"
	"math"
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
	//thread := uint64(1)

	wg := &sync.WaitGroup{}

	for i := uint64(0); i < thread; i++ {
		wg.Add(1)
		go prime(s, wg)
	}
	wg.Wait()
}

const (
	//bufferSize  = 41000000000
	//bufferLimit = 40999999990

	//*/
	bufferSize  = 200000000
	bufferLimit = 199990000
)

func prime(s *SafeNum, wg *sync.WaitGroup) {

	buf := make([]byte, 0, bufferSize)
	for num := s.Get(); num <= limit; num = s.Get() {

		for i := uint64(1); i != calcRange; i++ {
			num += 1
			if isPrime(num) {
				buf = append(append(buf, strconv.FormatUint(num, 10)...), '\n')

				if len(buf) > bufferLimit {
					_, err := os.Stdout.Write(buf)
					if err != nil {
						log.Fatal(err)
					}
					buf = make([]byte, 0, bufferSize)
				}
			}
			//fmt.Println(i, num)

		}

	}
	wg.Done()
	//os.Exit(0)
	//wg.Done()
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

	//zero := big.NewInt(0)
	//one := big.NewInt(1)
	//two := big.NewInt(2)
	//
	//if p.Cmp(two) == 0 {
	//	return true
	//}
	//
	//// p - 1 = 2^s * dに分解する
	//d := new(big.Int).Sub(p, one)
	//s := 0
	//for new(big.Int).And(d, one).Cmp(zero) == 0 {
	//	d.Rsh(d, 1)
	//	s++
	//}
	//
	//n := new(big.Int).Sub(p, one)
	//k := 20
	//rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	//
	//for i := 0; i < k; i++ {
	//	result := false
	//	// [1, n-1]からaをランダムに選ぶ
	//	a := new(big.Int).Rand(rnd, n)
	//	a.Add(a, one)
	//
	//	// a^{2^r * d} mod p != -1(= p - 1 = n)の比較
	//	tmp := new(big.Int).Exp(a, d, p)
	//	for r := 0; r < s; r++ {
	//		if tmp.Cmp(n) == 0 {
	//			result = true
	//			break
	//		}
	//		tmp.Exp(tmp, two, p)
	//	}
	//
	//	// a^d != 1 mod p の比較
	//	if !result && new(big.Int).Exp(a, d, p).Cmp(one) != 0 {
	//		return false
	//	}
	//}
	//
	//return true
}
