package channel

import (
	"testing"
)

var queue Queue

func init() {
	queue.Init()
}

func BenchmarkQueuePush1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		queue.Push("123456789012345678901234567890123456")
	}
}

func BenchmarkQueuePush2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		queue.Push("123456789012345678901234567890123456")
	}
}

func BenchmarkQueuePush3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		queue.Push("123456789012345678901234567890123456")
	}
}

func BenchmarkQueuePush4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		queue.Push("123456789012345678901234567890123456")
	}
}

func BenchmarkQueuePop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		queue.Pop()
	}
}
