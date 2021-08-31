package post

import (
	"database/sql"
	"log"
	"myreddit/pkg/user"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	return db, mock
}

func TestGetAll(t *testing.T) {
	db, mock := NewMock()
	tests := map[string]struct {
		repo  *PostRepo
		query string
		rows  *sqlmock.Rows
		want  []*Post
	}{
		"test1": {
			repo: NewPostRepo(db),
			query: `
				SELECT post_id, score, text, title, category, created_at, 
				 users.user_id, users.username
				FROM posts, users where posts.user_id = users.user_id;
			`,
			rows: sqlmock.NewRows(
				[]string{
					"post_id", "score", "text", "title", "category", "created_at", "user_id", "username",
				}).AddRow(
				"1", "0", "blabla", "lmao", "star-wars", "27.05.1996", "100500", "oleg",
			),
			want: []*Post{
				&Post{
					Score:     0,
					Id:        "1",
					Text:      "blabla",
					Title:     "lmao",
					Category:  "star-wars",
					CreatedAt: "27.05.1996",
					Author:    &user.User{Id: "100500", Login: "oleg"},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mock.ExpectQuery(tc.query).WillReturnRows(tc.rows)
			got, err := tc.repo.GetAll()
			assert.NoError(t, err)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %#v, got: %#v", tc.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestNewComment(t *testing.T) {
	db, mock := NewMock()
	tests := map[string]struct {
		repo        *PostRepo
		postId      string
		user        *user.User
		comment     *Comment
		insertQuery string
		getQuery    string
		rows        *sqlmock.Rows
		want        *Post
	}{
		"example1": {
			repo:    NewPostRepo(db),
			postId:  "0",
			user:    &user.User{Id: "120"},
			comment: &Comment{Id: "1", Body: "this is a body"},
			insertQuery: `
				INSERT INTO comments \(comment_id, post_id, user_id, text\)
				VALUES \(\$1, \$2, \$3, \$4\)
			`,
			getQuery: `
				SELECT post_id, score, text, title, category, created_at, posts.user_id, username 
				FROM posts, users 
				WHERE post_id = \$1 AND posts.user_id = users.user_id
			`,
			rows: sqlmock.NewRows(
				[]string{
					"post_id", "score", "text", "title", "category", "created_at", "user_id", "username",
				}).AddRow(
				"0", "0", "this is a body", "TIETLE", "star-wars", "28.05.2021", "120", "Oleg",
			),
			want: &Post{
				Id:        "0",
				Score:     0,
				Text:      "this is a body",
				Title:     "TIETLE",
				Category:  "star-wars",
				CreatedAt: "28.05.2021",
				Author:    &user.User{Id: "120", Login: "Oleg"},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mock.ExpectExec(tc.insertQuery).
				WithArgs(sqlmock.AnyArg(), tc.postId, tc.user.Id, tc.comment.Body).
				WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectQuery(tc.getQuery).WithArgs(tc.postId).WillReturnRows(tc.rows)
			got, err := tc.repo.NewComment(tc.postId, tc.user, tc.comment)
			assert.NoError(t, err)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %#v, got: %#v", tc.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteComment(t *testing.T) {
	bd, mock := NewMock()
	tests := map[string]struct {
		deleteQuery string
		getQuery    string
		repo        *PostRepo
		user        *user.User
		commentId   string
		postId      string
		rows        *sqlmock.Rows
		want        *Post
	}{
		"example1": {
			deleteQuery: `DELETE FROM comments WHERE comment_id = \$1`,
			getQuery: `
				SELECT post_id, score, text, title, category, created_at, posts.user_id, username 
				FROM posts, users 
				WHERE post_id = \$1 AND posts.user_id = users.user_id
			`,
			repo: NewPostRepo(bd),
			rows: sqlmock.NewRows(
				[]string{
					"post_id", "score", "text", "title", "category", "created_at", "user_id", "username",
				}).AddRow(
				"0", "0", "this is a body", "TIETLE", "star-wars", "28.05.2021", "120", "Oleg",
			),
			want: &Post{
				Id:        "0",
				Score:     0,
				Text:      "this is a body",
				Title:     "TIETLE",
				Category:  "star-wars",
				CreatedAt: "28.05.2021",
				Author:    &user.User{Id: "120", Login: "Oleg"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mock.ExpectExec(tc.deleteQuery).
				WithArgs(sqlmock.AnyArg()).
				WillReturnResult(sqlmock.NewResult(0, 0))
			mock.ExpectQuery(tc.getQuery).WithArgs(tc.postId).WillReturnRows(tc.rows)
			got, err := tc.repo.DeleteComment(tc.postId, tc.commentId, tc.user)
			assert.NoError(t, err)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %#v, got: %#v", tc.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeletePost(t *testing.T) {
	db, mock := NewMock()
	tests := map[string]struct {
		repo                *PostRepo
		postDeleteQuery     string
		commentsDeleteQuery string
		votesDeleteQuery    string
		getQuery            string
		user                *user.User
		postId              string
		rows                *sqlmock.Rows
		want                []*Post
	}{
		"example1": {
			repo:                NewPostRepo(db),
			postDeleteQuery:     `DELETE FROM posts WHERE post_id = \$1`,
			commentsDeleteQuery: `DELETE FROM comments WHERE post_id = \$1`,
			votesDeleteQuery:    `DELETE FROM votes WHERE post_id = \$1`,
			getQuery: `
				SELECT post_id, score, text, title, category, created_at, 
					users.user_id, users.username
				FROM posts, users where posts.user_id = users.user_id;
			`,
			user: &user.User{},
			rows: sqlmock.NewRows(
				[]string{
					"post_id", "score", "text", "title", "category", "created_at", "user_id", "username",
				}).
				AddRow(
					"1", "0", "blabla", "lmao", "star-wars", "27.05.1996", "100500", "oleg",
				).
				AddRow(
					"2", "1", "lsd", "lol", "harry-potter", "21.04.1292", "100501", "treville",
				),
			want: []*Post{
				&Post{
					Score: 0,
					Id:    "1", Text: "blabla", Title: "lmao", Category: "star-wars",
					CreatedAt: "27.05.1996",
					Author:    &user.User{Id: "100500", Login: "oleg"},
				},
				&Post{
					Score: 1,
					Id:    "2", Text: "lsd", Title: "lol", Category: "harry-potter",
					CreatedAt: "21.04.1292",
					Author:    &user.User{Id: "100501", Login: "treville"},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mock.ExpectExec(tc.commentsDeleteQuery).
				WithArgs(sqlmock.AnyArg()).
				WillReturnResult(sqlmock.NewResult(0, 0))
			mock.ExpectExec(tc.votesDeleteQuery).
				WithArgs(sqlmock.AnyArg()).
				WillReturnResult(sqlmock.NewResult(0, 0))
			mock.ExpectExec(tc.postDeleteQuery).
				WithArgs(sqlmock.AnyArg()).
				WillReturnResult(sqlmock.NewResult(0, 0))
			mock.ExpectQuery(tc.getQuery).
				WillReturnRows(tc.rows)
			got, err := tc.repo.DeletePost(tc.postId, tc.user)
			assert.NoError(t, err)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %#v, got: %#v", tc.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestNewPost(t *testing.T) {
	db, mock := NewMock()
	tests := map[string]struct {
		repo        *PostRepo
		insertQuery string
		user        *user.User
		post        *Post
		want        *Post
	}{
		"example1": {
			repo: NewPostRepo(db),
			insertQuery: `
				INSERT INTO posts \(post_id, text, title, category, user_id\) 
				VALUES \(\$1, \$2, \$3, \$4, \$5\)
			`,
			user: &user.User{Id: "42"},
			post: &Post{
				Id: "1", Text: "this is a text", Title: "revenge of the sith",
				Category: "star-wars",
			},
			want: &Post{
				Id: "1", Text: "this is a text", Title: "revenge of the sith",
				Category: "star-wars", Author: &user.User{Id: "42"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mock.ExpectExec(tc.insertQuery).
				WithArgs(sqlmock.AnyArg(), tc.post.Text,
					tc.post.Title, tc.post.Category, tc.user.Id).
				WillReturnResult(sqlmock.NewResult(1, 1))
			got, err := tc.repo.NewPost(tc.user, tc.post)
			assert.NoError(t, err)
			tc.want.Id = got.Id
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %#v, got: %#v", tc.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetAllByCategory(t *testing.T) {
	db, mock := NewMock()
	tests := map[string]struct {
		repo     *PostRepo
		category string
		query    string
		rows     *sqlmock.Rows
		want     []*Post
	}{
		"test1": {
			repo:     NewPostRepo(db),
			category: "star-wars",
			query: `
				SELECT post_id, score, text, title, category, created_at, 
				 users.user_id, users.username
				FROM posts, users 
				WHERE category = \$1 AND posts.user_id = users.user_id;
			`,
			rows: sqlmock.NewRows(
				[]string{
					"post_id", "score", "text", "title", "category", "created_at", "user_id", "username",
				}).AddRow(
				"1", "0", "blabla", "lmao", "star-wars", "27.05.1996", "100500", "oleg",
			),
			want: []*Post{
				&Post{
					Score:     0,
					Id:        "1",
					Text:      "blabla",
					Title:     "lmao",
					Category:  "star-wars",
					CreatedAt: "27.05.1996",
					Author:    &user.User{Id: "100500", Login: "oleg"},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mock.ExpectQuery(tc.query).WillReturnRows(tc.rows)
			got, err := tc.repo.GetAllByCategory(tc.category)
			assert.NoError(t, err)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %#v, got: %#v", tc.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetPostById(t *testing.T) {
	db, mock := NewMock()
	tests := map[string]struct {
		repo   *PostRepo
		postId string
		query  string
		rows   *sqlmock.Rows
		want   *Post
	}{
		"test1": {
			repo:   NewPostRepo(db),
			postId: "1",
			query: `
				SELECT post_id, score, text, 
					title, category, created_at, posts.user_id, username 
				FROM posts, users 
				WHERE post_id = \$1 AND posts.user_id = users.user_id
			`,
			rows: sqlmock.NewRows(
				[]string{
					"post_id", "score", "text", "title", "category", "created_at", "user_id", "username",
				}).AddRow(
				"1", "0", "blabla", "lmao", "star-wars", "27.05.1996", "100500", "oleg",
			),
			want: &Post{
				Score:     0,
				Id:        "1",
				Text:      "blabla",
				Title:     "lmao",
				Category:  "star-wars",
				CreatedAt: "27.05.1996",
				Author:    &user.User{Id: "100500", Login: "oleg"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mock.ExpectQuery(tc.query).WillReturnRows(tc.rows)
			got, err := tc.repo.GetPostById(tc.postId)
			assert.NoError(t, err)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %#v, got: %#v", tc.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestLike(t *testing.T) {
	db, mock := NewMock()
	tests := map[string]struct {
		repo             *PostRepo
		selectVotesQuery string
		insertVotesQuery string
		updatePostsQuery string
		selectPostsQuery string
		getVotesQuery    string
		user             *user.User
		postId           string
		want             *Post
		isLiked          bool
		voteRows         *sqlmock.Rows
		rows             *sqlmock.Rows
	}{
		"test1": {
			repo: NewPostRepo(db),
			selectVotesQuery: `
				SELECT \* FROM votes WHERE post_id = \$1 AND user_id = \$2
			`,
			insertVotesQuery: `
				INSERT INTO votes \(post_id, user_id\) VALUES \(\$1, \$2\)
			`,
			updatePostsQuery: `
				UPDATE posts SET score = score \+ 1 WHERE post_id = \$1
			`,
			selectPostsQuery: `
				SELECT post_id, score, text, title, category, created_at, posts.user_id, username 
				FROM posts, users 
				WHERE post_id = \$1 AND posts.user_id = users.user_id
			`,
			getVotesQuery: `
				SELECT users.user_id, users.username 
				FROM votes, users where votes.post_id = \$1
				AND votes.user_id = users.user_id
			`,
			user:   &user.User{Id: "100500"},
			postId: "1",
			voteRows: sqlmock.NewRows([]string{"user_id", "username"}).
				AddRow("1", "oleg").
				AddRow("2", "jack"),
			rows: sqlmock.NewRows(
				[]string{
					"post_id", "score", "text", "title", "category", "created_at", "user_id", "username",
				}).AddRow(
				"1", "1", "blabla", "lmao", "star-wars", "27.05.1996", "100500", "oleg",
			),
			want: &Post{
				Votes: []*user.User{
					{Id: "1", Login: "oleg"},
					{Id: "2", Login: "jack"},
				},
				Score: 1, Id: "1", Text: "blabla",
				Title: "lmao", Category: "star-wars", CreatedAt: "27.05.1996",
				Author: &user.User{Id: "100500", Login: "oleg"},
			},
			isLiked: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mock.ExpectQuery(tc.selectVotesQuery).
				WithArgs(tc.postId, tc.user.Id).
				WillReturnRows(sqlmock.NewRows([]string{"post_id", "user_id"}))
			mock.ExpectExec(tc.insertVotesQuery).
				WithArgs(tc.postId, tc.user.Id).
				WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec(tc.updatePostsQuery).
				WithArgs(tc.postId).
				WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectQuery(tc.selectPostsQuery).
				WithArgs(tc.postId).
				WillReturnRows(tc.rows)
			mock.ExpectQuery(tc.getVotesQuery).
				WithArgs(tc.postId).
				WillReturnRows(tc.voteRows)

			got, isLiked, err := tc.repo.Like(tc.user, tc.postId)
			log.Println(isLiked, err)
			assert.NoError(t, err)
			if tc.isLiked != isLiked {
				t.Fatalf("expected: %#v, got: %#v", tc.isLiked, isLiked)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %#v, got: %#v", tc.want, got)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
