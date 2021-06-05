package handler

import (
	"fmt"
	"log"
	"myreddit/pkg/post"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	PostRepo *post.PostRepo
}

func (h *PostHandler) GetAll(c *gin.Context) {
	posts, err := h.PostRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) GetAllByCategory(c *gin.Context) {
	category := c.Param("category")
	log.Println("cat: ", category)
	posts, err := h.PostRepo.GetAllByCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) GetById(c *gin.Context) {
	id := c.Param("id")
	post, err := h.PostRepo.GetPostById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	log.Println("About to return: ", post)
	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) NewComment(c *gin.Context) {
	user, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user is not found")
		return
	}

	comment := &post.Comment{}
	err = c.BindJSON(comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, "bad comment")
		return
	}

	log.Println("New comment: ", comment)

	postId := c.Param("id")
	post, err := h.PostRepo.NewComment(postId, user, comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) DeleteComment(c *gin.Context) {
	user, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, "no user")
		return
	}

	postId := c.Param("postId")
	commentId := c.Param("commentId")

	post, err := h.PostRepo.DeleteComment(postId, commentId, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	user, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, "no user")
		return
	}
	postId := c.Param("postId")
	posts, err := h.PostRepo.DeletePost(postId, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) NewPost(c *gin.Context) {
	u, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user is not found")
		return
	}

	p := &post.Post{}
	err = c.BindJSON(p)
	if err != nil {
		log.Println("BindJson: ", err)
		c.JSON(http.StatusInternalServerError, "oops")
		return
	}
	fmt.Println("\n\n\nPOST: ", p)
	post, err := h.PostRepo.NewPost(u, p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, post)
	log.Printf("Post: %v\n", post)
}

func (h *PostHandler) Like(c *gin.Context) {
	u, err := getUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, "user is not found")
		return
	}

	postId := c.Param("id")
	post, liked, err := h.PostRepo.Like(u, postId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"post": post,
		"like": liked,
	})
}
