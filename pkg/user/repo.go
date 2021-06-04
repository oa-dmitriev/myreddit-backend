package user

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/twinj/uuid"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "oleg"
	password = "12345"
	dbname   = "user"
)

type UserRepo struct {
	data map[string]*User
	Db   *sql.DB
}

func NewUserRepo() *UserRepo {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s sslmode=disable",
		host, port, user, password)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return &UserRepo{data: map[string]*User{}, Db: db}
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
