package post

import (
	"database/sql"
	"myreddit/pkg/user"

	"github.com/twinj/uuid"
)

type PostRepo struct {
	db *sql.DB
}

func NewPostRepo(db *sql.DB) *PostRepo {
	return &PostRepo{db}
}

func (repo *PostRepo) GetAll() ([]*Post, error) {
	sqlQuery := `
		SELECT post_id, score, text, title, category, created_at, 
		 users.user_id, users.username
		FROM posts, users where posts.user_id = users.user_id;
	`
	rows, err := repo.db.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := make([]*Post, 0)
	for rows.Next() {
		post := Post{}
		user := user.User{}
		err := rows.Scan(&post.Id, &post.Score, &post.Text, &post.Title,
			&post.Category, &post.CreatedAt, &user.Id, &user.Login)
		if err != nil {
			return nil, err
		}
		post.Author = &user
		post.Comments = getComments(repo.db, post.Id)
		post.Votes = getVotes(repo.db, post.Id)
		posts = append(posts, &post)
	}
	return posts, nil
}

func (repo *PostRepo) NewComment(postId string, u *user.User, comment *Comment) (*Post, error) {
	sqlStatement := `
		INSERT INTO comments (comment_id, post_id, user_id, text) 
		VALUES ($1, $2, $3, $4)
	`
	_, err := repo.db.Exec(sqlStatement,
		uuid.NewV4().String(), postId, u.Id, comment.Body)
	if err != nil {
		return nil, err
	}
	post, err := repo.GetPostById(postId)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (repo *PostRepo) DeleteComment(postId, commentId string, u *user.User) (*Post, error) {

	sqlStatement := "DELETE FROM comments WHERE comment_id = $1"
	_, err := repo.db.Exec(sqlStatement, commentId)
	if err != nil {
		return nil, err
	}
	post, err := repo.GetPostById(postId)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (repo *PostRepo) DeletePost(postId string, u *user.User) ([]*Post, error) {
	commentsDeleteQuery := "DELETE FROM comments WHERE post_id = $1"
	_, err := repo.db.Exec(commentsDeleteQuery, postId)
	if err != nil {
		return nil, err
	}
	votesDeleteQuery := "DELETE FROM votes WHERE post_id = $1"
	_, err = repo.db.Exec(votesDeleteQuery, postId)
	if err != nil {
		return nil, err
	}
	postDeleteQuery := "DELETE FROM posts WHERE post_id = $1"
	_, err = repo.db.Exec(postDeleteQuery, postId)
	if err != nil {
		return nil, err
	}
	return repo.GetAll()
}

func (repo *PostRepo) NewPost(u *user.User, p *Post) (*Post, error) {
	p.Id = uuid.NewV4().String()
	p.Author = u
	p.Score = 0
	sqlStatement := `
		INSERT INTO posts (post_id, text, title, category, user_id) 
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := repo.db.Exec(sqlStatement, p.Id, p.Text, p.Title, p.Category, u.Id)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (repo *PostRepo) GetAllByCategory(category string) ([]*Post, error) {
	sqlQuery := `
		SELECT post_id, score, text, title, category, created_at, 
		 users.user_id, users.username
		FROM posts, users 
		WHERE category = $1 AND posts.user_id = users.user_id;
	`
	rows, err := repo.db.Query(sqlQuery, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := make([]*Post, 0)
	for rows.Next() {
		post := Post{}
		user := user.User{}
		err := rows.Scan(&post.Id, &post.Score, &post.Text, &post.Title,
			&post.Category, &post.CreatedAt, &user.Id, &user.Login)
		if err != nil {
			return nil, err
		}
		post.Author = &user
		post.Comments = getComments(repo.db, post.Id)
		post.Votes = getVotes(repo.db, post.Id)
		posts = append(posts, &post)
	}
	return posts, nil
}

func (repo *PostRepo) GetPostById(postId string) (*Post, error) {
	var post Post
	author := user.User{}
	row := repo.db.QueryRow(`
		SELECT post_id, score, text, title, category, created_at, posts.user_id, username 
		FROM posts, users 
		WHERE post_id = $1 AND posts.user_id = users.user_id
	`, postId)
	err := row.Scan(&post.Id, &post.Score, &post.Text, &post.Title,
		&post.Category, &post.CreatedAt, &author.Id, &author.Login)
	if err != nil {
		return nil, err
	}
	post.Author = &author
	post.Comments = getComments(repo.db, post.Id)
	post.Votes = getVotes(repo.db, post.Id)
	return &post, nil
}

func (repo *PostRepo) Like(u *user.User, postId string) (*Post, bool, error) {
	row := repo.db.QueryRow(
		"SELECT * FROM votes WHERE post_id = $1 AND user_id = $2",
		postId, u.Id,
	)
	var pid, uid string
	if err := row.Scan(&pid, &uid); err == sql.ErrNoRows {
		_, err := repo.db.Exec(
			"INSERT INTO votes (post_id, user_id) VALUES ($1, $2)",
			postId, u.Id,
		)
		if err != nil {
			return nil, false, err
		}

		_, err = repo.db.Exec(
			"UPDATE posts SET score = score + 1 WHERE post_id = $1",
			postId,
		)
		if err != nil {
			return nil, false, err
		}
		p, err := repo.GetPostById(postId)
		if err != nil {
			return nil, false, err
		}
		p.Votes = getVotes(repo.db, postId)
		return p, true, nil
	}
	_, err := repo.db.Exec(`
		DELETE FROM votes 
		WHERE post_id = $1 AND user_id = $2
	`, postId, u.Id)
	if err != nil {
		return nil, false, err
	}
	_, err = repo.db.Exec(`
		UPDATE posts
		SET score = score - 1
		WHERE post_id = $1 
	`, postId)
	if err != nil {
		return nil, false, err
	}
	p, err := repo.GetPostById(postId)
	if err != nil {
		return nil, false, err
	}
	p.Votes = getVotes(repo.db, postId)
	return p, false, nil
}

func getVotes(db *sql.DB, postId string) []*user.User {
	sqlStatement := `
		SELECT users.user_id, users.username 
		FROM votes, users where votes.post_id = $1
		AND votes.user_id = users.user_id
	`
	rows, err := db.Query(sqlStatement, postId)
	if err != nil {
		return nil
	}

	users := make([]*user.User, 0)
	for rows.Next() {
		u := user.User{}
		err := rows.Scan(&u.Id, &u.Login)
		if err != nil {
			return nil
		}
		users = append(users, &u)
	}
	return users
}

func getComments(db *sql.DB, postId string) []*Comment {
	sqlStatement := `
		SELECT comment_id, users.user_id, username, text, created_at 
		FROM comments, users 
		WHERE post_id = $1 AND comments.user_id = users.user_id;
	`
	rows, err := db.Query(sqlStatement, postId)
	if err != nil {
		return nil
	}

	comments := make([]*Comment, 0)
	for rows.Next() {
		c := Comment{}
		u := user.User{}
		err := rows.Scan(&c.Id, &u.Id, &u.Login, &c.Body, &c.CreatedAt)
		if err != nil {
			return nil
		}
		c.Author = &u
		comments = append(comments, &c)
	}
	return comments
}
