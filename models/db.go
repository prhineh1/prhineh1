package models

import (
	"database/sql"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/mediocregopher/radix.v2/pool"
)

type Datastore interface {
	CreateUser(user *User) (string, error)
	VerifyLogin(ps, un string) error
	CreateSession(un, per string) (error, *http.Cookie)
	VerifySession(ck *http.Cookie) (string, error)
}

type DB struct {
	sql   *sql.DB
	cache *pool.Pool
}

func NewDB(dataSourceName string, test bool) (*DB, error) {
	var err error

	if test {
		db, err := sql.Open("postgres", dataSourceName)
		if err != nil {
			return nil, err
		}
		err = db.Ping()
		if err != nil {
			return nil, err
		}

		cache, err := pool.New("tcp", "localhost:3030", 10)
		if err != nil {
			return nil, err
		}

		return &DB{db, cache}, nil
	}

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	cache, err := pool.New("tcp", "localhost:6379", 10)
	if err != nil {
		return nil, err
	}

	return &DB{db, cache}, nil

}

// Only use for test database!!
func TearDown(db *DB) error {
	_, err := db.sql.Exec("DELETE FROM users")
	if err != nil {
		return err
	}
	return nil
}