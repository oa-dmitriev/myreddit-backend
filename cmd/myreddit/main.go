package main

import (
	"log"
	"myreddit/pkg/handler"
	"myreddit/pkg/middleware"
	"myreddit/pkg/post"
	"myreddit/pkg/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"token": "hola",
	})
}

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
	r.GET("/", test)
	r.POST("/api/register", userHandler.Register)
	r.POST("/api/login", userHandler.Login)

	r.GET("/api/posts", handler.GetAll)
	r.GET("/api/posts/:category", handler.GetAllByCategory)
	r.GET("/api/posts/:category/:id", handler.GetById)

	authorized := r.Group("/")
	authorized.Use(middleware.AuthRequired())
	log.Println("TELL ME WAT YOU NEED")
	{
		authorized.POST("/api/posts", handler.NewPost)
		authorized.POST("/api/posts/:id", handler.NewComment)
		authorized.DELETE("/api/posts/:postId/:commentId", handler.DeleteComment)
		authorized.DELETE("/api/posts/:postId", handler.DeletePost)
		authorized.GET("api/posts/:category/:id/like", handler.Like)
	}
	r.Run()
}
