package config

import (
	"log"
	"runtime"
)

// EnableProfilingMetrics - включение профилирования блокировок и мьютексов
func EnableProfilingMetrics() {
	// Включаем Block profile - показывает где горутины блокируются
	runtime.SetBlockProfileRate(1)

	// Включаем Mutex profile - показывает конфликты мьютексов
	runtime.SetMutexProfileFraction(1)

	log.Println("[PROFILE] Block and Mutex profiling enabled")
}

// GetMemStats - получение статистики памяти
func GetMemStats() runtime.MemStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m
}

// PrintMemStats - вывод статистики памяти в лог
func PrintMemStats() {
	m := GetMemStats()
	log.Printf("[MEMORY] Alloc: %v MB, TotalAlloc: %v MB, Sys: %v MB, NumGC: %v\n",
		m.Alloc/1024/1024,
		m.TotalAlloc/1024/1024,
		m.Sys/1024/1024,
		m.NumGC,
	)
}