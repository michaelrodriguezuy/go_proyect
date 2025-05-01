package user

import (
	"context"
	"log"

	"github.com/michaelrodriguezuy/go_proyect/internal/domain"
)

type (
	Service interface {
		Create(ctx context.Context, firstName, lastName string, age uint8) (*domain.User, error)
		GetAll(ctx context.Context) ([]domain.User, error)
		GetByID(ctx context.Context, id uint64) (*domain.User, error)
		Update(ctx context.Context, id uint64, firstName, lastName *string, age *uint8) error
	}

	service struct {
		log  *log.Logger
		repo Repository
	}
)

func NewService(logger *log.Logger, repo Repository) Service {
	return &service{
		log:  logger,
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, firstName, lastName string, age uint8) (*domain.User, error) {
	s.log.Println("Create")
	user := &domain.User{
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
	}
	savedUser, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return &savedUser, nil
}
func (s *service) GetAll(ctx context.Context) ([]domain.User, error) {
	s.log.Println("GetAll")
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *service) GetByID(ctx context.Context, id uint64) (*domain.User, error) {
	s.log.Println("GetByID")
	return s.repo.GetByID(ctx, id)
}

func (s *service) Update(ctx context.Context, id uint64, firstName, lastName *string, age *uint8) error {

	s.log.Println("Update")
	if err := s.repo.Update(ctx, id, firstName, lastName, age); err != nil {
		return err
	}
	return nil

}
