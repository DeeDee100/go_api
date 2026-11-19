package handlers

import (
	"net/http"
	"todo_api/internal/models"
	"todo_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func CreateUser(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input CreateUserInput
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		if len(input.Password) < 6 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters long!"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong validating password"})
		}

		user := &models.User{
			Email:    input.Email,
			Password: string(hashedPassword),
		}

		created, nil := repository.CreateUser(pool, user)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		}

		ctx.JSON(http.StatusCreated, created)
	}
}

func GetAllUser(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		users, err := repository.GetAllUser(pool)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		ctx.JSON(http.StatusOK, users)
	}
}
