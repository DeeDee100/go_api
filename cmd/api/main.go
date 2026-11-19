package main

import (
	"log"
	"todo_api/internal/config"
	postgres "todo_api/internal/database"
	"todo_api/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	var cfg *config.Config
	var err error

	cfg, err = config.Load()

	if err != nil {
		log.Fatal("Unable to load configuration:", err)
	}

	var pool *pgxpool.Pool

	pool, err = postgres.Connect(cfg.DatabaseURL)

	if err != nil {
		log.Fatal("Failed to connect to db:", err)
	}

	defer pool.Close()

	var router *gin.Engine = gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "TODO API is runnig!",
			"status":   "sucess",
			"database": "connected",
		})
	})

	router.POST("/todos", handlers.CreateTodoHandler(pool))
	router.GET("/todos", handlers.GetAllTodosHandler(pool))
	router.GET("/todos/:id", handlers.GetTodoByIDHandler(pool))
	router.PUT("/todos/:id", handlers.UpdateTodoHandler(pool))
	router.DELETE("/todos/:id", handlers.DeleteTodoHandler(pool))

	router.GET("/users", handlers.GetAllUser(pool))
	router.POST("/users", handlers.CreateUser(pool))

	// router.Run(":" + cfg.Port)
	router.Run(":8080")
}
