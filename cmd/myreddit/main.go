package main

import (
	"log"
	"myreddit/pkg/handler"
	"myreddit/pkg/middleware"
	"myreddit/pkg/post"
	"myreddit/pkg/user"
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

func main() {
	userRepo, err := user.NewUserRepo()
	if err != nil {
		log.Fatal("ERROR: ", err)
	}
	postRepo := post.NewPostRepo(userRepo.Db)
	userHandler := &handler.UserHandler{
		UserRepo: userRepo,
	}
	handler := &handler.PostHandler{
		PostRepo: postRepo,
	}

	r := gin.Default()
	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))
	r.POST("/api/register", userHandler.Register)
	r.POST("/api/login", userHandler.Login)
	r.GET("/api/category", handler.GetCategories)
	r.GET("/api/posts", handler.GetAll)
	r.GET("/api/posts/:category", handler.GetAllByCategory)
	r.GET("/api/posts/:category/:id", handler.GetById)

	authorized := r.Group("/")
	authorized.Use(middleware.AuthRequired())
	{
		authorized.POST("/api/category", handler.NewCategory)
		authorized.POST("/api/posts", handler.NewPost)
		authorized.POST("/api/posts/:id", handler.NewComment)
		authorized.DELETE("/api/posts/:postId/:commentId", handler.DeleteComment)
		authorized.DELETE("/api/posts/:postId", handler.DeletePost)
		authorized.GET("/api/posts/:category/:id/like", handler.Like)
	}
	r.Run()
}
