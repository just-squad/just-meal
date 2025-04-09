package http

import (
	"just-meal-api/internal/models"
	"just-meal-api/internal/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Server struct {
	Router   *gin.Engine
	dishRepo repositories.DishRepository
}

func NewWebApiServer(dishRepo repositories.DishRepository) *Server {
	server := &Server{
		Router:   gin.Default(),
		dishRepo: dishRepo,
	}
	server.setupMiddlewares()
	server.setupRoutes()
	return server
}

func (s *Server) setupMiddlewares() {
	s.Router.Use(gin.Logger())
	s.Router.Use(gin.Recovery())
}

func (s *Server) setupRoutes() {
	dishApi := s.Router.Group("/api/v1/dishes")
	{
		dishApi.POST("/", s.createDish)
		dishApi.GET("/", s.getDishesByType)
		dishApi.GET("/:id", s.getDish)
		dishApi.PUT("/:id", s.updateDish)
		dishApi.DELETE("/:id", s.deleteDish)
	}
}

func (s *Server) createDish(ctx *gin.Context) {
	var dish models.Dish
	if err := ctx.ShouldBindBodyWithJSON(&dish); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	dish.Id = uuid.New()
	if err := s.dishRepo.CreateDish(ctx.Request.Context(), &dish); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dish)
}

func (s *Server) getDishesByType(ctx *gin.Context) {
	mealType := models.MealType(ctx.Query("type"))

	if !isValidMealType(mealType) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid meal type"})
		return
	}

	dishes, err := s.dishRepo.GetDishesByType(ctx.Request.Context(), mealType)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dishes)
}

func (s *Server) getDish(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid Id format"})
		return
	}
	dish, err := s.dishRepo.GetDish(ctx.Request.Context(), id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Dish not found"})
	}

	ctx.JSON(http.StatusOK, dish)
}

func (s *Server) updateDish(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var dish models.Dish
	if err := ctx.ShouldBindJSON(&dish); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	dish.Id = id
	if err := s.dishRepo.UpdateDish(ctx.Request.Context(), id, &dish); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return
	}

	ctx.JSON(http.StatusOK, dish)
}

func (s *Server) deleteDish(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := s.dishRepo.DeleteDish(ctx.Request.Context(), id); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func isValidMealType(mt models.MealType) bool {
	switch mt {
	case models.Breakfast,
		models.Brunch,
		models.Dinner,
		models.Lunch,
		models.Snack,
		models.Supper:
		return true
	default:
		return false
	}
}
