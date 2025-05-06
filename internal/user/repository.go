package user

import (
	"context"

	"log"

	"database/sql"

	"github.com/michaelrodriguezuy/go_proyect/internal/domain"
)

/* type DB struct {
	Users     []domain.User
	MaxUserID uint64
} */

// los campos estan en minusculas porque van a ser campos privados
type (
	Repository interface {
		Create(ctx context.Context, user *domain.User) error
		GetAll(ctx context.Context) ([]domain.User, error)
		GetByID(ctx context.Context, id uint64) (*domain.User, error)
		Update(ctx context.Context, id uint64, firstName, lastName *string, age *uint8) error
		//Delete(id int) (domain.User, error)
	}

	repo struct {
		//db  DB //este es el mock
		db  *sql.DB
		log *log.Logger
	}
)

func NewRepository(db *sql.DB, logger *log.Logger) Repository {
	return &repo{
		db:  db,
		log: logger,
	}
}

func (r *repo) Create(ctx context.Context, user *domain.User) error {
	r.log.Println("Create repository")
	// Simulate a delay
	//time.Sleep(2 * time.Second)

	/* // Simulate a database insert
	user.ID = uint64(r.db.MaxUserID + 1)
	r.db.Users = append(r.db.Users, *user)
	r.db.MaxUserID++
	*/

	//inyecto el usuario directamente a la base de datos
	sqlQ := `INSERT INTO users (first_name, last_name, age) VALUES (?, ?, ?)`

	res, err := r.db.Exec(sqlQ, user.FirstName, user.LastName, user.Age)

	if err != nil {
		r.log.Println("Error inserting user:", err.Error())
		return err
	}

	id, err := res.LastInsertId()

	if err != nil {
		r.log.Println("Error getting last insert id:", err.Error())
	}
	user.ID = uint64(id)

	r.log.Printf("User created with ID: ", id)
	return nil
}
func (r *repo) GetAll(ctx context.Context) ([]domain.User, error) {
	r.log.Println("GetAll")
	// Simulate a delay
	// time.Sleep(2 * time.Second)
	//return r.db.Users, nil

	var users []domain.User
	sqlQ := `SELECT id, first_name, last_name, age FROM users`
	rows, err := r.db.Query(sqlQ)
	if err != nil {
		r.log.Println("Error querying users:", err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Age); err != nil {
			r.log.Println("Error scanning user:", err.Error())
			return nil, err
		}
		users = append(users, user)
	}
	r.log.Println("Users found: ", len(users))
	return users, nil
}

func (r *repo) GetByID(ctx context.Context, id uint64) (*domain.User, error) {
	r.log.Println("GetByID")
	// Simulate a delay
	// time.Sleep(2 * time.Second)

	//esto me retorna el indice del elemento que busco, y lo busco por el id
	//sino encuentra el elemento, devuelve -1

	//de esta forma acceso a un campo dentro de una estructura

	/*
		index := slices.IndexFunc(r.db.Users, func(user domain.User) bool {
			return user.ID == id
		})

		if index < 0 {
			return nil, ErrUserNotFound{id}
		}

		return &r.db.Users[index], nil
	*/

	return nil, nil

}

func (r *repo) Update(ctx context.Context, id uint64, firstName, lastName *string, age *uint8) error {
	r.log.Println("Update")

	//como el get trabaja con punteros, no es necesario hacer un cast, los cambios se hacen directamente en memoria
	/*
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
	*/

	return nil

}
