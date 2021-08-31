# My reddit clone

## Description 

Reddit is a social network where you can post, comment and like posts published by others

### Main endpoints:

- POST   /api/register - registration 
- POST   /api/login - login
- GET    /api/posts - get all posts
- GET    /api/posts/category - get all posts of the chosen category
- GET    /api/posts/category/id - get the post by id 

### Only authenticated users get to use the following endpoints 

- POST   /api/category - create a new category
- POST   /api/posts - create a new post
- POST   /api/posts/id - leave a comment with id 
- DELETE /api/posts/postId/commentId - delete the comment with commentId
- DELETE /api/posts/postId - delete the post with postId
- GET    /api/posts/category/id/like - leave a like to the post with id
