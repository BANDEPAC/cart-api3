package app

import (
	"cart-api/internal/config"
	"cart-api/internal/db/postgres"
	"cart-api/internal/service"
	handler "cart-api/internal/transport/http"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run initializes the application by loading configuration, connecting to the database,
// setting up service and repository layers, and starting the HTTP server with defined routes.
func Run() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	db, err := postgres.Connect(cfg)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	cartRepo := postgres.NewCartRepository(db)
	cartitemRepo := postgres.NewCartItemRepository(db)
	cartService := service.NewCartService(cartRepo)
	cartitemService := service.NewCartItemRepository(cartitemRepo)
	cartHandler := handler.NewCartHandler(cartService, cartitemService)

	router := http.NewServeMux()

	router.Handle("POST /carts", http.HandlerFunc(cartHandler.CreateCart))
	router.Handle("GET /carts/{id}", http.HandlerFunc(cartHandler.ViewCart))
	router.Handle("POST /carts/{id}/items", http.HandlerFunc(cartHandler.AddToCart))
	router.Handle("DELETE /carts/{id}/items/{item_id}", http.HandlerFunc(cartHandler.RemoveFromCart))

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server is running on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v", cfg.ServerPort, err)
		}
	}()

	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")

}
