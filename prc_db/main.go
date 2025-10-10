package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// fallback — прямой DSN в коде (только для учебного стенда!)
		dsn = "postgres://postgres:@localhost:5433/prc_db"
	}

	db, err := openDB(dsn)
	if err != nil {
		log.Fatalf("openDB error: %v", err)
	}
	defer db.Close()

	repo := NewRepo(db)

	// 1) Вставим пару задач
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	titles := []string{"Сделать ПЗ №5", "Купить кофе", "Проверить отчёты"}
	for _, title := range titles {
		id, err := repo.CreateTask(ctx, title)
		if err != nil {
			log.Fatalf("CreateTask error: %v", err)
		}
		log.Printf("Inserted task id=%d (%s)", id, title)
	}

	// 2) Прочитаем список задач
	ctxList, cancelList := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelList()

	tasks, err := repo.ListTasks(ctxList)
	if err != nil {
		log.Fatalf("ListTasks error: %v", err)
	}

	// 3) Напечатаем
	fmt.Println("=== Tasks ===")
	for _, t := range tasks {
		fmt.Printf("#%d | %-24s | done=%-5v | %s\n",
			t.ID, t.Title, t.Done, t.CreatedAt.Format(time.RFC3339))
	}

	ctxList, cancelList = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelList()

	tasks, err = repo.ListDone(ctxList, true)
	if err != nil {
		log.Fatalf("ListDone error: %v", err)
	}

	// 3) Напечатаем
	fmt.Println("=== Tasks ===")
	for _, t := range tasks {
		fmt.Printf("#%d | %-24s | done=%-5v | %s\n",
			t.ID, t.Title, t.Done, t.CreatedAt.Format(time.RFC3339))
	}

	ctxList, cancelList = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelList()

	t, err := repo.FindByID(ctxList, 3)
	if err != nil {
		log.Fatalf("FindByID error: %v", err)
	}

	// 3) Напечатаем
	fmt.Println("=== Task ===")
	fmt.Printf("#%d | %-24s | done=%-5v | %s\n",
		t.ID, t.Title, t.Done, t.CreatedAt.Format(time.RFC3339))

	ctxList, cancelList = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelList()

	titles = []string{"Проверка 1", "Проверка 2", "Проверка 3"}
	dones := []bool{true, true, true}
	err = repo.CreateMany(ctxList, titles, dones)
	if err != nil {
		log.Fatalf("CreateMany error: %v", err)
	}
	log.Printf("Tasks have been inserted")
}
