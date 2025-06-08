package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/kenwoo9y/todo-api-go/api/internal/config"
	"github.com/kenwoo9y/todo-api-go/api/internal/db"
	"github.com/kenwoo9y/todo-api-go/api/internal/handler"
	"github.com/kenwoo9y/todo-api-go/api/internal/repository"
	"github.com/kenwoo9y/todo-api-go/api/internal/server"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.New()
	if err != nil {
		return err
	}

	// DB接続
	database, err := db.NewDB()
	if err != nil {
		return err
	}
	defer database.Close()

	// リポジトリの初期化
	userRepo := repository.NewUserRepository(database, cfg)
	taskRepo := repository.NewTaskRepository(database, cfg)

	// ハンドラーの初期化
	userHandler := handler.NewUserHandler(userRepo)
	taskHandler := handler.NewTaskHandler(taskRepo)

	// サーバーの設定
	s := server.SetupServer(userHandler, taskHandler)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}

	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)

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
