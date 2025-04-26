package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/michaelrodriguezuy/go_proyect/internal/user"
	"github.com/michaelrodriguezuy/go_proyect/pkg/transport"
)

func NewUserHTTPServer(ctx context.Context, router *http.ServeMux, endpoints user.Endpoints) {
	router.HandleFunc("/users/", UserServer(ctx, endpoints)) //para usar el handle con el router interno de GO, necesitamos el request, que nos lo da la funcion UserServer
}

func UserServer(ctx context.Context, endpoints user.Endpoints) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		url := r.URL.Path
		log.Println(r.Method, ": ", url)

		path, pathSize := transport.Clean(url)
		if pathSize < 3 || pathSize > 4 {
			return //si el path no es correcto, no hace nada
		}

		//aca obtengo el id del usuario
		params := make(map[string]string)
		if pathSize == 4 && path[2] != "" {
			params["userId"] = path[2]
		}

		//agrego el id al contexto
		ctx = context.WithValue(ctx, "params", params)

		tran := transport.NewTransport(w, r, ctx) //con este transport lo que hacemos es generar los middleware

		switch r.Method {
		case http.MethodGet:
			switch pathSize {
			case 3: //es un getall
				tran.Server(
					transport.Endpoint(endpoints.GetAll),
					decodeGetAllUser,
					encodeResponse,
					encodeError)
			case 4: //es un get por id
				log.Println("Get por id")
				tran.Server(
					nil,
					decodeGetUser,
					encodeResponse,
					encodeError)
			}
			return
		case http.MethodPost:
			switch pathSize {
			case 3: //es un create
				tran.Server(
					transport.Endpoint(endpoints.Create),
					decodeCreateUser,
					encodeResponse,
					encodeError)
			}
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func decodeGetUser(ctx context.Context, r *http.Request) (any, error) {
	params := ctx.Value("params").(map[string]string)

	fmt.Println("params: ", params)
	fmt.Println("userId: ", params["userId"])
	return nil, fmt.Errorf("not implemented")
}

func decodeGetAllUser(ctx context.Context, r *http.Request) (any, error) {
	return nil, nil
}

func decodeCreateUser(ctx context.Context, r *http.Request) (any, error) {
	var req user.CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { //decodifica el body
		return nil, fmt.Errorf("error decoding request: %v", err.Error())
	}
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp any) error {

	data, err := json.Marshal(resp) //marshall convierte una estructura a un json

	if err != nil {
		return err
	}

	status := http.StatusOK

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, `{"status": %d, "data": "%s"}`, status, data)

	// if err := json.NewEncoder(w).Encode(data); err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }
	return nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	status := http.StatusInternalServerError
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, `{"status": %d, "data": "%s"}`, status, err.Error())
}
