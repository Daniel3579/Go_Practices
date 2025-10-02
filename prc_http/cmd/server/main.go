package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"example.com/prc_http/internal/api"
	"example.com/prc_http/internal/storage"
)

func main() {
	store := storage.NewMemoryStore()
	h := api.NewHandlers(store)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		api.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Коллекция
	mux.HandleFunc("GET /tasks", h.ListTasks)
	mux.HandleFunc("POST /tasks", h.CreateTask)
	// Элемент
	mux.HandleFunc("GET /tasks/", h.GetTask)
	mux.HandleFunc("PATCH /tasks/", h.PatchTask)
	mux.HandleFunc("DELETE /tasks/", h.DeleteTask)

	// Сначала оборачиваем mux в CORS
	handler := api.CORS(mux)

	// Затем подключаем логирование
	handler = api.Logging(handler)

	// Получаем значение порта из переменной окружения PRC_HTTP_PORT
	portStr := os.Getenv("PRC_HTTP_PORT")
	if portStr == "" {
		portStr = "8080"
	}

	// Преобразуем строку в целое число
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port number: %s\n", portStr)
		return
	}

	// Создаем HTTP-сервер
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Println("listening on", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Обработка сигналов завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Создаем контекст с таймаутом для завершения работы
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Корректное завершение работы сервера
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exited gracefully")
}
