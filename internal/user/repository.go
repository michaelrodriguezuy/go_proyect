package user

import (
	"context"
	"log"

	"github.com/michaelrodriguezuy/go_proyect/internal/domain"
)

type DB struct {
	Users     []domain.User
	MaxUserID uint64
}

// los campos estan en minusculas porque van a ser campos privados
type (
	Repository interface {
		GetAll(ctx context.Context) ([]domain.User, error)
		//GetByID(id int) (domain.User, error)
		Create(ctx context.Context, user *domain.User) (domain.User, error)
		//Update(user domain.User) (domain.User, error)
		//Delete(id int) (domain.User, error)
	}

	repo struct {
		db  DB
		log *log.Logger
	}
)

func NewRepository(db DB, logger *log.Logger) Repository {
	return &repo{
		db:  db,
		log: logger,
	}
}

func (r *repo) Create(ctx context.Context, user *domain.User) (domain.User, error) {
	r.log.Println("Create")
	// Simulate a delay
	//time.Sleep(2 * time.Second)
	user.ID = int(r.db.MaxUserID + 1)
	r.db.Users = append(r.db.Users, *user)
	r.db.MaxUserID++
	r.log.Printf("User created: %+v\n", user)
	return *user, nil
}
func (r *repo) GetAll(ctx context.Context) ([]domain.User, error) {
	r.log.Println("GetAll")
	// Simulate a delay
	// time.Sleep(2 * time.Second)
	return r.db.Users, nil
}
