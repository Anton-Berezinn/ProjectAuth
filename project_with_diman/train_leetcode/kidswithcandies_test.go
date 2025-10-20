package main

import "testing"

func BenchmarkKidsWithCandies(b *testing.B) {
	value := []int{2, 3, 5, 1, 3}
	extra := 3
	for i := 0; i < b.N; i++ {
		kidsWithCandies(value, extra)
	}
}

func BenchmarkKidsWithCandies2(b *testing.B) {
	value := []int{2, 3, 5, 1, 3}
	extra := 3
	for i := 0; i < b.N; i++ {
		kidsWithCandies(value, extra)
	}
}
