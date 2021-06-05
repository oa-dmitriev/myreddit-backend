package main

import (
	"log"
	"myreddit/pkg/handler"
	"myreddit/pkg/middleware"
	"myreddit/pkg/post"
	"myreddit/pkg/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://myreddit-frontend.herokuapp.com"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	r.POST("/api/register", userHandler.Register)
	r.POST("/api/login", userHandler.Login)

	r.GET("/api/posts", handler.GetAll)
	r.GET("/api/posts/:category", handler.GetAllByCategory)
	r.GET("/api/posts/:category/:id", handler.GetById)

	authorized := r.Group("/")
	authorized.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://myreddit-frontend.herokuapp.com"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	authorized.Use(middleware.AuthRequired())
	{
		authorized.POST("/api/posts", handler.NewPost)
		authorized.POST("/api/posts/:id", handler.NewComment)
		authorized.DELETE("/api/posts/:postId/:commentId", handler.DeleteComment)
		authorized.DELETE("/api/posts/:postId", handler.DeletePost)
		authorized.GET("api/posts/:category/:id/like", handler.Like)
	}
	r.Run()
}
