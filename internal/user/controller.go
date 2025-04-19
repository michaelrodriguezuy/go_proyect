package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

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

func NewEndpoint(ctx context.Context, service Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			decode := json.NewDecoder(r.Body) //decodifica el body
			var req CreateReq

			if err := decode.Decode(&req); err != nil { //decodifica el body y lo guarda en la variable user
				MessageResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			PostUser(ctx, service, w, req)
		case http.MethodGet:
			GetAllUser(ctx, service, w)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func GetAllUser(ctx context.Context, service Service, w http.ResponseWriter) {
	service.GetAll(ctx)
	users, err := service.GetAll(ctx)
	if err != nil {
		MessageResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	DataResponse(w, http.StatusOK, users)
}

func PostUser(ctx context.Context, service Service, w http.ResponseWriter, data any) { //lo mismo, aca es una interfaz, pero en las versiones nuevas se usa any

	req := data.(CreateReq) //convierte el dato a un tipo User CASTEO

	if req.FirstName == "" || req.LastName == "" {
		MessageResponse(w, http.StatusBadRequest, "Faltan datos")
		return
	}
	if req.Age < 18 {
		MessageResponse(w, http.StatusBadRequest, "El usuario no es mayor de edad")
		return
	}

	user, err := service.Create(ctx, req.FirstName, req.LastName, req.Age)
	if err != nil {
		MessageResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	DataResponse(w, http.StatusCreated, user) 

	fmt.Println("Usuario creado:", user)
}



func DataResponse(w http.ResponseWriter, status int, data any) { //en este caso users es una entidad del tipo interfaz(any en las ultimas versiones), la idea con esto es que le pueda mandar cualquier entidad y la tome

	value, err := json.Marshal(data) //marshall convierte una estructura a un json

	if err != nil {
		MessageResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "data": "%s"}`, status, value)

	// if err := json.NewEncoder(w).Encode(data); err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }
}

func MessageResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "message": "%s"}`, status, message)
}
