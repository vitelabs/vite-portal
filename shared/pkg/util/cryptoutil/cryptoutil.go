package cryptoutil

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	mathrand "math/rand"
)

// Source: https://blog.gopheracademy.com/advent-2017/a-tale-of-two-rands/
func NewCryptoSeededSource() mathrand.Source {
	var seed int64
	binary.Read(cryptorand.Reader, binary.BigEndian, &seed)
	return mathrand.NewSource(seed)
}

// RangeRandomInt returns a random range of numbers as a slice within a range
func RangeRandomInt(min int, max int, n int) []int {
	if max-min <= 0 || n == 0 {
		return make([]int, 0)
	}
	arr := make([]int, n)
	random := mathrand.New(NewCryptoSeededSource())
	for r := 0; r < n; r++ {
		arr[r] = random.Intn(max-min) + min
	}
	return arr
}

// UniqueRandomInt returns a random range of unique numbers as a slice
func UniqueRandomInt(max int, n int) []int {
	if max < n || n == 0 {
		return make([]int, 0)
	}
	arr := make([]int, n)
	if max == n {
		for i := 0; i < n; i++ {
			arr[i] = i
		}
		return arr
	}
	random := mathrand.New(NewCryptoSeededSource())
	m := make(map[int]bool)
	for i, curr := 0, 0; ; i++ {
		r := random.Intn(max)

		// check if random number is already in the map -> continue
		if _, ok := m[r]; ok {
			continue
		}

		// add the random number to the map
		m[r] = true
		// add the random number to the array
		arr[curr] = r

		curr++
		// end loop if reached n random numbers
		if curr == n {
			break
		}
	}
	return arr
}
