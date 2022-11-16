package client

import "testing"

func BenchmarkGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(&GetOptions{Env: "test"})
	}
}
