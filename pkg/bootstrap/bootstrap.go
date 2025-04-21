package bootstrap

import (
	"log"
	"os"

	"github.com/michaelrodriguezuy/go_proyect/internal/domain"
	"github.com/michaelrodriguezuy/go_proyect/internal/user"
)

func NewLogger() *log.Logger {
	return log.New(os.Stdout, "user: ", log.LstdFlags|log.Lshortfile)
}

func NewDB() user.DB {

	return user.DB{
		Users: []domain.User{
			{ID: 1, FirstName: "John", LastName: "Doe", Age: 30},
			{ID: 2, FirstName: "Jane", LastName: "Smith", Age: 25},
		},
		MaxUserID: 2,
	}
}
