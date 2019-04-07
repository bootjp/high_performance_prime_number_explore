package main

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	limit     = uint64(math.MaxUint32)
	calcRange = 1000
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
	thread := uint64(runtime.NumCPU() * 2)
	//thread := uint64(1)

	wg := &sync.WaitGroup{}

	for i := uint64(0); i < thread; i++ {
		wg.Add(1)
		go prime(s, wg)
	}
	wg.Wait()
}

func prime(s *SafeNum, wg *sync.WaitGroup) {
	for num := s.Get(); num <= limit; num = s.Get() {

		for i := uint64(1); i != calcRange; i++ {
			num += 1
			if isPrime(big.NewInt(int64(num))) {
				fmt.Println(num)
			}
			//fmt.Println(i, num)

		}

	}
	os.Exit(0)
	//wg.Done()
}

func isPrime(p *big.Int) bool {
	zero := big.NewInt(0)
	one := big.NewInt(1)
	two := big.NewInt(2)

	if p.Cmp(two) == 0 {
		return true
	}

	// p - 1 = 2^s * dに分解する
	d := new(big.Int).Sub(p, one)
	s := 0
	for new(big.Int).And(d, one).Cmp(zero) == 0 {
		d.Rsh(d, 1)
		s++
	}

	n := new(big.Int).Sub(p, one)
	k := 20
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < k; i++ {
		result := false
		// [1, n-1]からaをランダムに選ぶ
		a := new(big.Int).Rand(rnd, n)
		a.Add(a, one)

		// a^{2^r * d} mod p != -1(= p - 1 = n)の比較
		tmp := new(big.Int).Exp(a, d, p)
		for r := 0; r < s; r++ {
			if tmp.Cmp(n) == 0 {
				result = true
				break
			}
			tmp.Exp(tmp, two, p)
		}

		// a^d != 1 mod p の比較
		if !result && new(big.Int).Exp(a, d, p).Cmp(one) != 0 {
			return false
		}
	}

	return true
}
