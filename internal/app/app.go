package app

import (
	"cart-api/internal/config"
	"cart-api/internal/db/postgres"
	"cart-api/internal/service"
	handler "cart-api/internal/transport/http"
	"log"
	"net/http"
	"strings"
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

	router.HandleFunc("/carts", cartHandler.CreateCart)
	router.HandleFunc("/carts/", func(w http.ResponseWriter, r *http.Request) {
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

		log.Println("Parsed path:", pathParts)

		if len(pathParts) < 2 || pathParts[0] != "carts" {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		switch {
		case len(pathParts) == 2 && r.Method == http.MethodGet:
			cartHandler.ViewCart(w, r)
		case len(pathParts) == 3 && pathParts[2] == "items" && r.Method == http.MethodPost:
			cartHandler.AddToCart(w, r)
		case len(pathParts) == 4 && pathParts[2] == "items" && r.Method == http.MethodDelete:
			cartHandler.RemoveFromCart(w, r)

		default:
			http.Error(w, "Invalid request method or path", http.StatusMethodNotAllowed)
		}
	})

	log.Printf("Server is running on port %s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, router))
}
