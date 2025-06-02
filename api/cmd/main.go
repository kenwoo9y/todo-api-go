package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/kenwoo9y/todo-api-go/api/internal/config"
)

func run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.New()
	if err != nil {
		return err
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}

	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)

	mux := http.NewServeMux()

	// タスク関連のエンドポイント
	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// タスク一覧の取得
			fmt.Fprintf(w, "Get tasks")
		case http.MethodPost:
			// タスクの作成
			fmt.Fprintf(w, "Create task")
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// 特定のタスクのエンドポイント
	mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/tasks/")
		if id == "" {
			http.Error(w, "Task ID is required", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			// 特定のタスクの取得
			fmt.Fprintf(w, "Get task: %s", id)
		case http.MethodPut:
			// タスクの更新
			fmt.Fprintf(w, "Update task: %s", id)
		case http.MethodDelete:
			// タスクの削除
			fmt.Fprintf(w, "Delete task: %s", id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// ユーザー関連のエンドポイント
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// ユーザー一覧の取得
			fmt.Fprintf(w, "Get users")
		case http.MethodPost:
			// ユーザーの作成
			fmt.Fprintf(w, "Create user")
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// 特定のユーザーのエンドポイント
	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/users/")
		if id == "" {
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			// 特定のユーザーの取得
			fmt.Fprintf(w, "Get user: %s", id)
		case http.MethodPut:
			// ユーザーの更新
			fmt.Fprintf(w, "Update user: %s", id)
		case http.MethodDelete:
			// ユーザーの削除
			fmt.Fprintf(w, "Delete user: %s", id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	s := &http.Server{
		Handler: mux,
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.Serve(l); err != nil &&
			err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
		}
	}()

	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}

	wg.Wait()
	return nil
}
