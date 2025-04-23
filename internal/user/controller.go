package user

import (
	"context"
	"errors"
	"log"
)

type (
	Controller func(ctx context.Context, req any) (any, error)

	Endpoints struct {
		//estas funciones siguen el patron de funcion declarado en la linea 10
		Create Controller
		GetAll Controller
	}

	CreateReq struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name,omitempty"`
		Age       uint8  `json:"age"`
	}
)

// este metodo me devuelve una estructura de endpoints, linea 12
func NewEndpoint(ctx context.Context, service Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(service),
		GetAll: makeGetAllEndpoint(service),
		//cuando estas funciones se ejecutan, no hace un return de las funciones de la linea 45 y 34
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
		log.Println("Usuario creado:", user)
		return user, nil
	}
}
