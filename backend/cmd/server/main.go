package main

import (
	"database/sql"
	"fmt"
	"log"

	"bathroom-admin/internal/config"
	"bathroom-admin/internal/handler"
	"bathroom-admin/internal/middleware"
	"bathroom-admin/internal/model"
	"bathroom-admin/pkg/jwt"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	jwtMgr := jwt.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpireHours)
	userRepo := model.NewUserRepository(db)
	authHandler := handler.NewAuthHandler(userRepo, jwtMgr)
	authMiddleware := middleware.NewAuthMiddleware(jwtMgr)

	categoryRepo := model.NewCategoryRepository(db)
	categoryHandler := handler.NewCategoryHandler(categoryRepo)

	productRepo := model.NewProductRepository(db)
	productHandler := handler.NewProductHandler(productRepo)

	if cfg.App.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authMiddleware.RequireAuth(), authHandler.Logout)
		}

		categories := v1.Group("/categories")
		{
			categories.GET("", categoryHandler.List)
			categories.POST("", authMiddleware.RequireAuth(), categoryHandler.Create)
		}

		products := v1.Group("/products")
		{
			products.GET("", productHandler.List)
			products.GET("/:id", productHandler.Detail)
			products.POST("", authMiddleware.RequireAuth(), productHandler.Create)
			products.PUT("/:id", authMiddleware.RequireAuth(), productHandler.Update)
			products.DELETE("/:id", authMiddleware.RequireAuth(), productHandler.Delete)
		}
	}

	addr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)
	log.Printf("Server running on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
