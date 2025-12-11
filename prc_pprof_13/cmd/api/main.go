package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	"example.com/prc_pprof/internal/config"
	"example.com/prc_pprof/internal/work"
)

func main() {
	config.EnableProfilingMetrics()

	// Эндпоинт для CPU-нагрузки (рекурсивный Фибоначчи)
	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		defer work.TimeIt("Fib(38) - исходная версия")()
		n := 38
		res := work.Fib(n)
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(fmt.Sprintf("%d\n", res)))
	})

	// Эндпоинт для оптимизированной версии
	http.HandleFunc("/work-fast", func(w http.ResponseWriter, r *http.Request) {
		defer work.TimeIt("FibFast(38) - оптимизированная версия")()
		n := 38
		res := work.FibFast(n)
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(fmt.Sprintf("%d\n", res)))
	})

	// Эндпоинт для тестирования горутин и памяти
	http.HandleFunc("/allocate", func(w http.ResponseWriter, r *http.Request) {
		defer work.TimeIt("allocate - тест аллокаций")()
		data := work.AllocateMemory(1000000)
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(fmt.Sprintf("allocated %d items\n", len(data))))
	})

	// Эндпоинт для оптимизированной версии
	http.HandleFunc("/allocate-fast", func(w http.ResponseWriter, r *http.Request) {
		defer work.TimeIt("allocateFast - тест аллокаций")()
		data := work.FastAllocateMemory(1000000)
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(fmt.Sprintf("allocated %d items\n", len(data))))
	})

	// Эндпоинт для тестирования горутин
	http.HandleFunc("/goroutines", func(w http.ResponseWriter, r *http.Request) {
		defer work.TimeIt("goroutines - запуск горутин")()
		count := work.SpawnGoroutines(100)
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(fmt.Sprintf("spawned %d goroutines\n", count)))
	})

	// Хеалтчек
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("ok\n"))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
