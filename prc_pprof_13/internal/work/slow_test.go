package work

import (
	"testing"
)

// BenchmarkFib - бенчмарк для рекурсивного Фибоначчи
func BenchmarkFib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Fib(30)
	}
}

// BenchmarkFibFast - бенчмарк для итеративного Фибоначчи
func BenchmarkFibFast(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = FibFast(30)
	}
}

// BenchmarkFibMemo - бенчмарк для мемоизированного Фибоначчи
func BenchmarkFibMemo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = FibMemo(30)
	}
}

// BenchmarkAllocateMemory - бенчмарк для аллокации памяти
func BenchmarkAllocateMemory(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = AllocateMemory(1000000)
	}
}

// BenchmarkSpawnGoroutines - бенчмарк для создания горутин
func BenchmarkSpawnGoroutines(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = SpawnGoroutines(100)
	}
}

// BenchmarkSpawnGoroutinesPool - бенчмарк для pool pattern
func BenchmarkSpawnGoroutinesPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = SpawnGoroutinesPool(100, 10)
	}
}

// BenchmarkFib38 - отдельный бенчмарк для Fib(38) (более тяжелый)
func BenchmarkFib38(b *testing.B) {
	b.Run("Recursive", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Fib(30)
		}
	})
	b.Run("Iterative", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = FibFast(30)
		}
	})
	b.Run("Memoized", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = FibMemo(30)
		}
	})
}

// BenchmarkAllocations - бенчмарк со сфокусом на аллокации
func BenchmarkAllocations(b *testing.B) {
	b.Run("Regular", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = AllocateMemory(100000)
		}
	})
	b.Run("Fast", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = FastAllocateMemory(100000)
		}
	})
}

func TestDummy(t *testing.T) {
	// Пустой тест - просто чтобы Go знал, что тесты есть
}