package user

import (
	"context"
	"errors"
	"fmt"
	

	"github.com/michaelrodriguezuy/go_proyect2/response"
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
			return nil, response.BadRequest(ErrFirstNameRequired.Error())
		}
		if req.LastName == "" {
			return nil, response.BadRequest(ErrLastNameRequired.Error())
		}
		if req.Age < 18 {
			return nil, response.BadRequest(ErrAgeMinor.Error())
		}

		if err := service.Create(ctx, req.FirstName, req.LastName, req.Age); err != nil {
			return nil, response.InternalServerError(err.Error())
		}
		
		return response.Created("success", nil), nil
	}
}
func makeGetAllEndpoint(service Service) Controller {

	return func(ctx context.Context, req any) (any, error) {
		users, err := service.GetAll(ctx)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}
		return response.OK("success ", users), nil
	}
}

func makeGetByIDEndpoint(service Service) Controller {

	return func(ctx context.Context, request any) (any, error) {

		req := request.(GetReq) //casteo el dato a un tipo User
		fmt.Println("req: ", req)

		user, err := service.GetByID(ctx, req.ID)
		if err != nil {
			if errors.As(err, &ErrUserNotFound{}) {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success ", user), nil
	}
}

func makeUpdateEndpoint(service Service) Controller {
	return func(ctx context.Context, request any) (any, error) {

		req := request.(UpdateReq)

		if req.FirstName != nil && *req.FirstName == "" {
			return nil, response.BadRequest(ErrFirstNameRequired.Error())
		}
		if req.LastName != nil && *req.LastName == "" {
			return nil, response.BadRequest(ErrLastNameRequired.Error())
		}
		if req.Age != nil && *req.Age < 18 {
			return nil, response.BadRequest(ErrAgeMinor.Error())
		}

		if err := service.Update(ctx, req.ID, req.FirstName, req.LastName, req.Age); err != nil {

			if errors.As(err, &ErrUserNotFound{}) {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", nil), nil
	}
}
