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
	"github.com/Nerfi/instaClone/internal/handlers/middlewares"
	postsHndlr "github.com/Nerfi/instaClone/internal/handlers/posts"
	authRepo "github.com/Nerfi/instaClone/internal/repository/authRepo"
	postsRepo "github.com/Nerfi/instaClone/internal/repository/posts"
	authSrv "github.com/Nerfi/instaClone/internal/services/auth"
	postsSrv "github.com/Nerfi/instaClone/internal/services/posts"
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

	postsRepository := postsRepo.NewPostsRepo(cfg.DB)
	// 3 init service
	authSrv := authSrv.NewAuthSrv(authRepo, cfg.Auth)
	postsService := postsSrv.NewPostsSrv(postsRepository)

	//4 init handlers
	authHandlers := authHndlr.NewAuthHanlders(authSrv)
	postsHandlers := postsHndlr.NewPostsHanlders(postsService)
	// 5 register routes
	mux := http.NewServeMux()
	mux.HandleFunc("/create", authHandlers.CreateUser)
	mux.HandleFunc("/login", authHandlers.LoginUser)
	mux.HandleFunc("/refresh", authHandlers.RefreshToken)
	mux.Handle("/logout", middlewares.AuthMiddleware(http.HandlerFunc(authHandlers.LogoutUser)))
	mux.Handle("/forgot-password", middlewares.AuthMiddleware(http.HandlerFunc(authHandlers.ForgotPassword)))

	// healthCheck endpoint
	hc := health.NewHealthCheck()
	mux.HandleFunc("/health", hc.Check)

	// aplicar el CSRF
	//TODO descomentar csrf logic para probarla cuando tengamos un endpoint valido
	//csrfMiddleware := security.NewCSRF([]byte(config.Envs.CSRF_SECRET_KEY), false)(mux)

	// chain middleware
	chain := middlewares.ChainMiddleware(middlewares.AuthMiddleware, middlewares.OwnerOnlyMiddleware)
	// rutas protegidas
	mux.Handle("/profile/{id}", chain(http.HandlerFunc(authHandlers.Profile)))

	// POSTS routes
	mux.HandleFunc("/posts", postsHandlers.GetPosts)
	// tiene que estar protegido
	mux.Handle("/post/create", middlewares.AuthMiddleware(http.HandlerFunc(postsHandlers.CreatePost)))
	//delete endpoint, should be secure
	mux.Handle("/posts/{id}", middlewares.AuthMiddleware(http.HandlerFunc(postsHandlers.DeletePost)))

	// get single post, not sure if this one needs auth
	mux.Handle("/post/{id}", middlewares.AuthMiddleware(http.HandlerFunc(postsHandlers.GetPost)))

	//securing headers in all requests coming through this router

	secureMux := middlewares.SecurityHeaders(mux)

	// custom server config

	s := &http.Server{
		Addr:    ":8081",
		Handler: secureMux,
		//Handler:        csrfMiddleware,
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
