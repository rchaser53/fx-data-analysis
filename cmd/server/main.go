package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rchaser53/fx-data-analysis/internal/database"
	"github.com/rchaser53/fx-data-analysis/internal/handler"
)

func main() {
	// Initialize database
	db, err := database.NewDB("./fx_trades.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize schema
	if err := db.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	// Initialize handler
	h := handler.NewHandler(db)

	// Setup router
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// API routes
	api := r.Group("/api/v1")
	{
		trades := api.Group("/trades")
		{
			trades.POST("", h.CreateTrade)
			trades.GET("", h.GetAllTrades)
			trades.GET("/:id", h.GetTrade)
			trades.PUT("/:id", h.UpdateTrade)
			trades.DELETE("/:id", h.DeleteTrade)
		}
	}

	// Start server
	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
