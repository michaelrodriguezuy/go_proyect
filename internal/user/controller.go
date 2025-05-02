package user

import (
	"context"
	"fmt"
	"log"
)

type (
	Controller func(ctx context.Context, req any) (any, error)

	Endpoints struct {
		//estas funciones siguen el patron de funcion declarado en la linea 10
		Create  Controller
		GetAll  Controller
		GetByID Controller
		Update  Controller
	}

	GetReq struct {
		ID uint64
	}

	CreateReq struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name,omitempty"`
		Age       uint8  `json:"age"`
	}

	//con el puntero lo que le digo es que si viene vacio es un string vacio, y si no viene el valor del campo es nulo, aca no lo actualiza
	UpdateReq struct {
		ID        uint64  `json:"id"`
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name,omitempty"`
		Age       *uint8  `json:"age"`
	}
)

// este metodo me devuelve una estructura de endpoints, linea 12
func NewEndpoint(ctx context.Context, service Service) Endpoints {
	return Endpoints{
		//cuando estas funciones se ejecutan, no hace un return de las funciones de la linea 45 y 34
		Create: makeCreateEndpoint(service),
		GetAll: makeGetAllEndpoint(service),

		GetByID: makeGetByIDEndpoint(service),
		Update:  makeUpdateEndpoint(service),
	}
}

func makeCreateEndpoint(service Service) Controller {
	return func(ctx context.Context, request any) (any, error) {

		req := request.(CreateReq) //convierte el dato a un tipo User CASTEO

		if req.FirstName == "" {
			return nil, ErrFirstNameRequired
		}
		if req.LastName == "" {
			return nil, ErrLastNameRequired
		}
		if req.Age < 18 {
			return nil, ErrAgeMinor
		}

		user, err := service.Create(ctx, req.FirstName, req.LastName, req.Age)
		if err != nil {
			return nil, err
		}
		log.Println("Usuario creado:", user)
		return user, nil
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

func makeGetByIDEndpoint(service Service) Controller {

	return func(ctx context.Context, request any) (any, error) {

		req := request.(GetReq) //casteo el dato a un tipo User
		fmt.Println("req: ", req)

		users, err := service.GetByID(ctx, req.ID)
		if err != nil {
			return nil, err
		}

		return users, nil
	}
}

func makeUpdateEndpoint(service Service) Controller {
	return func(ctx context.Context, request any) (any, error) {

		req := request.(UpdateReq)

		if req.FirstName != nil && *req.FirstName == "" {
			return nil, ErrFirstNameRequired
		}
		if req.LastName != nil && *req.LastName == "" {
			return nil, ErrLastNameRequired
		}
		if req.Age != nil && *req.Age < 18 {
			return nil, ErrAgeMinor
		}

		if err := service.Update(ctx, req.ID, req.FirstName, req.LastName, req.Age); err != nil {

			return nil, err
		}

		return nil, nil
	}
}
