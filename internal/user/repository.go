package user

import (
	"context"

	"log"

	"slices"

	"github.com/michaelrodriguezuy/go_proyect/internal/domain"
)

type DB struct {
	Users     []domain.User
	MaxUserID uint64
}

// los campos estan en minusculas porque van a ser campos privados
type (
	Repository interface {
		Create(ctx context.Context, user *domain.User) (domain.User, error)
		GetAll(ctx context.Context) ([]domain.User, error)
		GetByID(ctx context.Context, id uint64) (*domain.User, error)
		Update(ctx context.Context, id uint64, firstName, lastName *string, age *uint8) error
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
	user.ID = uint64(r.db.MaxUserID + 1)
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

func (r *repo) GetByID(ctx context.Context, id uint64) (*domain.User, error) {
	r.log.Println("GetByID")
	// Simulate a delay
	// time.Sleep(2 * time.Second)

	//esto me retorna el indice del elemento que busco, y lo busco por el id
	//sino encuentra el elemento, devuelve -1

	//de esta forma acceso a un campo dentro de una estructura
	index := slices.IndexFunc(r.db.Users, func(user domain.User) bool {
		return user.ID == id
	})

	if index < 0 {
		return nil, ErrUserNotFound{id}
	}

	return &r.db.Users[index], nil

}

func (r *repo) Update(ctx context.Context, id uint64, firstName, lastName *string, age *uint8) error {
	r.log.Println("Update")

	//como el get trabaja con punteros, no es necesario hacer un cast, los cambios se hacen directamente en memoria
	user, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if firstName != nil {
		user.FirstName = *firstName
	}
	if lastName != nil {
		user.LastName = *lastName
	}
	if age != nil {
		user.Age = *age
	}

	return nil

}
