package user

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/twinj/uuid"
)

const (
	host     = "ec2-54-228-9-90.eu-west-1.compute.amazonaws.com"
	port     = 5432
	user     = "hqkdqdwfzzmknz"
	password = "5ea140a1ab2205f26672db6020f5297b1f696c66db9202b6841dc763d1726394"
	dbname   = "d1ec1tpipgs10c"
)

type UserRepo struct {
	data map[string]*User
	Db   *sql.DB
}

func NewUserRepo() (*UserRepo, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			user_id uuid NOT NULL,
			username varchar(60) NOT NULL,
			password varchar(60) NOT NULL,
			created_on date NOT NULL DEFAULT NOW(),
			CONSTRAINT PK_users PRIMARY KEY ( user_id )
		)
	`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			post_id    uuid NOT NULL,
			score      int NOT NULL,
			text       text NOT NULL,
			title      varchar(255) NOT NULL,
			category   varchar(255) NOT NULL,
			created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
			user_id    uuid NOT NULL,
			CONSTRAINT PK_posts PRIMARY KEY ( post_id ),
			CONSTRAINT FK_23 FOREIGN KEY ( user_id ) REFERENCES users ( user_id )
		 )
	`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS comments (
			comment_id uuid NOT NULL,
			post_id    uuid NOT NULL,
			user_id    uuid NOT NULL,
			text       text NOT NULL,
			created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT PK_comments PRIMARY KEY ( comment_id ),
			CONSTRAINT FK_39 FOREIGN KEY ( post_id ) REFERENCES posts ( post_id ),
			CONSTRAINT FK_55 FOREIGN KEY ( user_id ) REFERENCES users ( user_id )
		)
	`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS votes (
			post_id uuid NOT NULL,
			user_id uuid NOT NULL,
			CONSTRAINT FK_60 FOREIGN KEY ( post_id ) REFERENCES posts ( post_id ),
			CONSTRAINT FK_65 FOREIGN KEY ( user_id ) REFERENCES users ( user_id )
		)
	`)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS categories (
			name				text NOT NULL,
			description text NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}
	return &UserRepo{data: map[string]*User{}, Db: db}, nil
}

func (repo *UserRepo) Register(u *User) error {
	u.Id = uuid.NewV4().String()

	sqlStatement := `
		INSERT INTO users (user_id, username, password) 
		VALUES ($1, $2, $3)
	`
	_, err := repo.Db.Exec(sqlStatement, u.Id, u.Login, u.Password)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepo) Login(user *User) error {
	var u User
	sqlStatement := `
		SELECT user_id, username FROM users where username = $1 AND password = $2
	`
	err := repo.Db.QueryRow(sqlStatement, user.Login, user.Password).
		Scan(&u.Id, &u.Login)
	if err != nil {
		return fmt.Errorf("login or password are not valid")
	}
	user.Id = u.Id
	return nil
}
