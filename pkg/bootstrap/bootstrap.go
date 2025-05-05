package bootstrap

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" //lo importo pero no lo uso, por eso el guion bajo
	_ "github.com/michaelrodriguezuy/go_proyect/internal/domain"
	_ "github.com/michaelrodriguezuy/go_proyect/internal/user"
)

func NewLogger() *log.Logger {
	return log.New(os.Stdout, "user: ", log.LstdFlags|log.Lshortfile)
}

func NewDBConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3307)/go_users")

	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

/*
func NewDB() user.DB {

	return user.DB{
		Users: []domain.User{
			{ID: 1, FirstName: "John", LastName: "Doe", Age: 30},
			{ID: 2, FirstName: "Jane", LastName: "Smith", Age: 25},
		},
		MaxUserID: 2,
	}
}
*/
