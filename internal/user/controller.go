package user

import (
	"context"
	"errors"
	"fmt"
)

type (
	Controller func(ctx context.Context, req any) (any, error)

	Endpoints struct {
		Create Controller
		GetAll Controller
	}

	CreateReq struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name,omitempty"`
		Age       uint8  `json:"age"`
	}
)

func NewEndpoint(ctx context.Context, service Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(service),
		GetAll: makeGetAllEndpoint(service),
	}
}

func makeGetAllEndpoint(service Service) Controller {

	return func(ctx context.Context, req any) (any, error) {
		users, err := service.GetAll(ctx)
		if err != nil {
			return nil, err
		}
		return users, nil
	}
}

func makeCreateEndpoint(service Service) Controller {
	return func(ctx context.Context, request any) (any, error) {

		req := request.(CreateReq) //convierte el dato a un tipo User CASTEO

		if req.FirstName == "" || req.LastName == "" {
			return nil, errors.New("Faltan datos")
		}
		if req.Age < 18 {
			return nil, errors.New("El usuario no es mayor de edad")
		}

		user, err := service.Create(ctx, req.FirstName, req.LastName, req.Age)
		if err != nil {
			return nil, err
		}
		fmt.Println("Usuario creado:", user)
		return user, nil
	}
}
