package work

import "sync"

// FastAllocateMemory - оптимизированная аллокация с pre-allocated слайсом
func FastAllocateMemory(size int) []int64 {
	data := make([]int64, size)
	for i := 0; i < size; i++ {
		data[i] = int64(i * i)
	}
	return data
}

// SpawnGoroutinesPool - версия с ограничением горутин (pool pattern)
func SpawnGoroutinesPool(count int, poolSize int) int {
	done := make(chan bool, poolSize)
	sem := make(chan bool, poolSize)

	go func() {
		for i := 0; i < count; i++ {
			go func() {
				sem <- true
				defer func() { <-sem }()

				sum := 0
				for j := 0; j < 1000000; j++ {
					sum += j
				}
				done <- true
			}()
		}
	}()

	for i := 0; i < count; i++ {
		<-done
	}
	return count
}

// BatchProcessing - обработка данных батчами (пример оптимизации)
func BatchProcessing(data []int, batchSize int) int64 {
	var result int64
	var mu sync.Mutex

	batchCount := (len(data) + batchSize - 1) / batchSize

	done := make(chan bool, batchCount)

	for i := 0; i < batchCount; i++ {
		start := i * batchSize
		end := start + batchSize
		if end > len(data) {
			end = len(data)
		}

		go func(batch []int) {
			sum := int64(0)
			for _, v := range batch {
				sum += int64(v)
			}
			mu.Lock()
			result += sum
			mu.Unlock()
			done <- true
		}(data[start:end])
	}

	for i := 0; i < batchCount; i++ {
		<-done
	}

	return result
}
