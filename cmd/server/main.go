package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ahmed/auth-service/internal/config"
	"github.com/ahmed/auth-service/pkg/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	log.Println("Configuration loaded successfully")
	// 2. Connect to Database
	db, err := database.Connect(cfg.DB.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close() // Ensure connection is closed when main exits
	log.Printf("Connected to database '%s' successfully", cfg.DB.Name)
	// 3. Set up Router using Chi
	r := chi.NewRouter()
	r.Use(middleware.Logger)    // Log HTTP requests
	r.Use(middleware.Recoverer) // Prevent server crashes if a panic occurs
	// Dummy health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	// 4. Configure HTTP Server timeouts
	// Never run http.ListenAndServe() without timeouts in production!
	// Unfinished requests could hang open forever and run out of memory.
	server := &http.Server{
		Addr:         cfg.Server.Addr(),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	// 5. Start HTTP Server asynchronously
	go func() {
		log.Printf("Server starting on %s", cfg.Server.Addr())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server listen failed: %v", err)
		}
	}()

}
