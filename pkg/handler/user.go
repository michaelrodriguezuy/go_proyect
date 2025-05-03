package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/michaelrodriguezuy/go_proyect/internal/user"
	"github.com/michaelrodriguezuy/go_proyect/pkg/transport"
	"github.com/michaelrodriguezuy/go_proyect2/response"
)

func NewUserHTTPServer(ctx context.Context, router *http.ServeMux, endpoints user.Endpoints) {
	router.HandleFunc("/users/", UserServer(ctx, endpoints)) //para usar el handle con el router interno de GO, necesitamos el request, que nos lo da la funcion UserServer
}

func UserServer(ctx context.Context, endpoints user.Endpoints) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		url := r.URL.Path
		log.Println(r.Method, ": ", url)

		path, pathSize := transport.Clean(url)

		//aca obtengo el id del usuario
		params := make(map[string]string)
		if pathSize == 4 && path[2] != "" {
			params["userId"] = path[2]
		}

		//agrego el id al contexto, y con este transport lo que hacemos es generar los middleware
		tran := transport.NewTransport(w, r, context.WithValue(ctx, "params", params))

		var end user.Controller
		var deco func(ctx context.Context, r *http.Request) (any, error)

		switch r.Method {
		case http.MethodGet:
			switch pathSize {
			case 3: //es un getall
				end = endpoints.GetAll
				deco = decodeGetAllUser

			case 4: //es un get por id
				end = endpoints.GetByID
				deco = decodeGetUser

			}
		case http.MethodPost:
			switch pathSize {
			case 3: //es un create
				end = endpoints.Create
				deco = decodeCreateUser
			}

		case http.MethodPatch:
			switch pathSize {
			case 4: //es un update
				end = endpoints.Update
				deco = decodeUpdateUser
			}
		}

		if end != nil && deco != nil {
			tran.Server(
				transport.Endpoint(end),
				deco,
				encodeResponse,
				encodeError)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	}
}

func decodeUpdateUser(ctx context.Context, r *http.Request) (any, error) {

	var req user.UpdateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil { //decodifica el body
		return nil, fmt.Errorf("error decoding request: %v", err.Error())
	}

	params := ctx.Value("params").(map[string]string)

	id, err := strconv.ParseUint(params["userId"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing userId: %v", err.Error())
	}

	req.ID = id

	return req, nil

}

func decodeGetUser(ctx context.Context, r *http.Request) (any, error) {
	params := ctx.Value("params").(map[string]string)

	id, err := strconv.ParseUint(params["userId"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing userId: %v", err.Error())
	}

	fmt.Println("params: ", params)
	fmt.Println("userId: ", params["userId"])
	return user.GetReq{ID: id}, nil
}

func decodeGetAllUser(ctx context.Context, r *http.Request) (any, error) {
	fmt.Println("entro")
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

	//agrego el manejo del pkg response

	r := resp.(response.Response) //casteo el dato a un tipo Response
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())

	return json.NewEncoder(w).Encode(resp) //encodeo el response
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response) //casteo el dato a un tipo Response
	w.WriteHeader(resp.StatusCode())

	_ = json.NewEncoder(w).Encode(resp) //encodeo el response
}
