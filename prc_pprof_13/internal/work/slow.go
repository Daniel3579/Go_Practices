package work

// Fib - неоптимальный рекурсивный Фибоначчи (демонстрация CPU-нагрузки)
func Fib(n int) int {
	if n < 2 {
		return n
	}
	return Fib(n-1) + Fib(n-2)
}

// FibMemo - мемоизированный Фибоначчи (промежуточная оптимизация)
func FibMemo(n int) int {
	memo := make(map[int]int)
	var fib func(int) int
	fib = func(x int) int {
		if x < 2 {
			return x
		}
		if val, ok := memo[x]; ok {
			return val
		}
		memo[x] = fib(x-1) + fib(x-2)
		return memo[x]
	}
	return fib(n)
}

// FibFast - оптимизированный итеративный Фибоначчи (O(n) время, O(1) память)
func FibFast(n int) int {
	if n < 2 {
		return n
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

// AllocateMemory - функция для теста аллокаций памяти
func AllocateMemory(size int) []int64 {
    data := []int64{}  // Пустой слайс, будет расти
    for i := 0; i < size; i++ {
        data = append(data, int64(i*i))  // Каждый append может переаллоцировать
    }
    return data
}



// SpawnGoroutines - функция для теста горутин
func SpawnGoroutines(count int) int {
	done := make(chan bool, count)
	for i := 0; i < count; i++ {
		go func() {
			// Имитируем какую-то работу
			sum := 0
			for j := 0; j < 1000000; j++ {
				sum += j
			}
			done <- true
		}()
	}
	for i := 0; i < count; i++ {
		<-done
	}
	return count
}