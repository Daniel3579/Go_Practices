package work

import (
	"log"
	"time"
)

// TimeIt - простой декоратор для измерения времени выполнения функции
func TimeIt(name string) func() {
	start := time.Now()
	return func() {
		elapsed := time.Since(start)
		log.Printf("[TIMER] %s took %v (%d ms)\n", name, elapsed, elapsed.Milliseconds())
	}
}

// TimeItNano - версия таймера с наносекундной точностью
func TimeItNano(name string) func() {
	start := time.Now()
	return func() {
		elapsed := time.Since(start)
		log.Printf("[TIMER] %s took %d ns\n", name, elapsed.Nanoseconds())
	}
}

// Stopwatch - структура для более сложного профилирования
type Stopwatch struct {
	name  string
	start time.Time
}

// NewStopwatch - создание нового таймера
func NewStopwatch(name string) *Stopwatch {
	return &Stopwatch{
		name:  name,
		start: time.Now(),
	}
}

// Stop - остановка таймера и вывод результата
func (s *Stopwatch) Stop() time.Duration {
	elapsed := time.Since(s.start)
	log.Printf("[STOPWATCH] %s took %v\n", s.name, elapsed)
	return elapsed
}

// Lap - промежуточное время (без остановки таймера)
func (s *Stopwatch) Lap(label string) time.Duration {
	elapsed := time.Since(s.start)
	log.Printf("[LAP] %s - %s: %v\n", s.name, label, elapsed)
	return elapsed
}
