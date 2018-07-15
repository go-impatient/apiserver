package util

import (
	"math/rand"
	"time"
)

// RandInts returns a random integer array with the specified from, to and size.
func RandInts(from, to, size int) []int {
	if to-from < size {
		size = to - from
	}

	var slice []int
	for i := from; i < to; i++ {
		slice = append(slice, i)
	}

	var ret []int
	for i := 0; i < size; i++ {
		idx := rand.Intn(len(slice))
		ret = append(ret, slice[idx])
		slice = append(slice[:idx], slice[idx+1:]...)
	}

	return ret
}

// RandString returns a fixed length random string.
func RandString(n int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	time.Sleep(time.Nanosecond)

	letter := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}

	return string(b)
}
