package handler

import (
	"net/http"
	"context"
	"encoding/json"
	"fmt"
	
	"github.com/michaelrodriguezuy/go_proyect/internal/user"
	"github.com/michaelrodriguezuy/go_proyect/pkg/transport"
)


func NewUserHTTPServer(ctx context.Context, router *http.ServeMux, endpoints user.Endpoints) {	
	router.HandleFunc("/users", UserServer(ctx, endpoints))
}

func UserServer(ctx context.Context, endpoints user.Endpoints) func(w http.ResponseWriter, r *http.Request) {	
		return func(w http.ResponseWriter, r *http.Request) {

			transport := transport.NewTransport(w, r, ctx) 

			switch r.Method {
			case http.MethodGet:
				transport.Server(
					transport.Endpoint(endpoints.GetAll),
					decodeGetAllUser,
					encodeResponse,
				)

					
			case http.MethodPost:
				decode := json.NewDecoder(r.Body) //decodifica el body
				var req CreateReq
		
				if err := decode.Decode(&req); err != nil { //decodifica el body y lo guarda en la variable user
					MessageResponse(w, http.StatusBadRequest, err.Error())
					return
				}
				PostUser(ctx, service, w, req)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}
}

func decodeGetAllUser(ctx context.Context, r *http.Request) (any, error) {
	return nil, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp any) error {

	data, err := json.Marshal(resp) //marshall convierte una estructura a un json
	
	if err != nil {		
		return err
	}

	status := http.StatusOK
	
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "data": "%s"}`, status, data)
	
	// if err := json.NewEncoder(w).Encode(data); err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }
}


func DataResponse(w http.ResponseWriter, status int, data any) { //en este caso users es una entidad del tipo interfaz(any en las ultimas versiones), la idea con esto es que le pueda mandar cualquier entidad y la tome

}

func MessageResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "message": "%s"}`, status, message)
}