package main

import (
	"fmt"
	"math"
)

func main() {
	var isPrime bool
	count := 2
	fmt.Printf("1\t2\n")
	for i := uint64(3); i < math.MaxUint64; i += 2 {
		isPrime = true
		for j := uint64(3); j*j <= i; j += 2 {
			if i%j == 0 {
				isPrime = false
				break
			}
		}
		if isPrime {
			fmt.Printf("%v\t%v\n", count, i)
			count++
		}
	}
}
