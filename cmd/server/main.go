package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Nerfi/instaClone/internal/config"
	authHndlr "github.com/Nerfi/instaClone/internal/handlers/auth"
	health "github.com/Nerfi/instaClone/internal/handlers/healthCheck"
	authRepo "github.com/Nerfi/instaClone/internal/repository/authRepo"
	authSrv "github.com/Nerfi/instaClone/internal/services/auth"
	"github.com/joho/godotenv"
)

// init load the enviroment variables from a .env file and parses them
func init() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file")
		panic(err)
	}

	// cargamos las envs para que puedan ser usadas en todo el proyecto
	_, err = config.ParseEnvs()
	if err != nil {
		panic(err)
	}
}

func main() {
	// 1 init config
	cfg := config.NewAppConfig()
	// 2 init repo
	authRepo := authRepo.NewAuthRepo(cfg.DB)
	// 3 init service
	authSrv := authSrv.NewAuthSrv(authRepo, cfg.Auth)

	//4 init handlers
	authHandlers := authHndlr.NewAuthHanlders(authSrv)
	// 5 register routes
	mux := http.NewServeMux()
	mux.HandleFunc("/create", authHandlers.CreateUser)
	mux.HandleFunc("/login", authHandlers.LoginUser)

	// healthCheck endpoint
	hc := health.NewHealthCheck()
	mux.HandleFunc("/health", hc.Check)

	// custom server config

	s := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// graceful shutdown

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}

	}()
	log.Printf("server listening at %s", s.Addr)

	// wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// shutdown server gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server force to shutdown:", err)
	}
	log.Println("Server gracefully shutdown")
}
