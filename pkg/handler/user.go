package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/michaelrodriguezuy/go_proyect/internal/user"
	"github.com/michaelrodriguezuy/go_proyect/pkg/transport"
	"github.com/michaelrodriguezuy/go_proyect2/response"
)

func NewUserHTTPServer(endpoints user.Endpoints) http.Handler {
	r := gin.Default()

	r.POST("/users", transport.GinServer(
		transport.Endpoint(endpoints.Create),
		decodeCreateUser,
		encodeResponse,
		encodeError,
	))
	r.GET("/users", transport.GinServer(
		transport.Endpoint(endpoints.GetAll),
		decodeGetAllUser,
		encodeResponse,
		encodeError,
	))
	r.GET("/users/:id", transport.GinServer(
		transport.Endpoint(endpoints.GetByID),
		decodeGetUser,
		encodeResponse,
		encodeError,
	))
	r.PATCH("/users/:id", transport.GinServer(
		transport.Endpoint(endpoints.Update),
		decodeUpdateUser,
		encodeResponse,
		encodeError,
	))

	return r
}

// func UserServer( endpoints user.Endpoints) func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		url := r.URL.Path
// 		log.Println(r.Method, ": ", url)

// 		path, pathSize := transport.Clean(url)

// 		//aca obtengo el id del usuario
// 		params := make(map[string]string)
// 		if pathSize == 4 && path[2] != "" {
// 			params["userId"] = path[2]
// 		}

// 		//envio el token en el param para poder hacer en cada endpoint la verificacion, quizas quiera que alguno de los endpoints sea publico
// 		//verifico si el token es correcto
// 		token := r.Header.Get("Authorization")
// 		params["token"] = token

// 		//agrego el id al contexto, y con este transport lo que hacemos es generar los middleware
// 		tran := transport.NewTransport(w, r, context.WithValue(ctx, "params", params))

// 		var end user.Controller
// 		var deco func(ctx context.Context, r *http.Request) (any, error)

// 		switch r.Method {
// 		case http.MethodGet:
// 			switch pathSize {
// 			case 3: //es un getall
// 				end = endpoints.GetAll
// 				deco = decodeGetAllUser

// 			case 4: //es un get por id
// 				end = endpoints.GetByID
// 				deco = decodeGetUser

// 			}
// 		case http.MethodPost:
// 			switch pathSize {
// 			case 3: //es un create
// 				end = endpoints.Create
// 				deco = decodeCreateUser
// 			}

// 		case http.MethodPatch:
// 			switch pathSize {
// 			case 4: //es un update
// 				end = endpoints.Update
// 				deco = decodeUpdateUser
// 			}
// 		}

// 		if end != nil && deco != nil {
// 			tran.Server(
// 				transport.Endpoint(end),
// 				deco,
// 				encodeResponse,
// 				encodeError)
// 		} else {
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		}

// 	}
// }

func decodeUpdateUser(ctx *gin.Context) (any, error) {

	if err := tokenVerify(ctx.Request.Header.Get("Authorization")); err != nil {
		return nil, response.Unauthorized(err.Error())
	}

	var req user.UpdateReq
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil { //decodifica el body
		return nil, response.BadRequest(fmt.Sprintf("error decoding request: %v", err.Error()))
	}

	id, err := strconv.ParseUint(ctx.Params.ByName("id"), 10, 64)
	if err != nil {
		return nil, response.BadRequest(fmt.Sprintf("error parsing userId: %v", err.Error()))
	}

	req.ID = id

	return req, nil

}

func decodeGetUser(ctx *gin.Context) (any, error) {

	id, err := strconv.ParseUint(ctx.Params.ByName("id"), 10, 64)
	if err != nil {
		return nil, response.BadRequest(fmt.Sprintf("error parsing userId: %v", err.Error()))
	}

	return user.GetReq{ID: id}, nil
}

func decodeGetAllUser(ctx *gin.Context) (any, error) {
	if err := tokenVerify(ctx.Request.Header.Get("Authorization")); err != nil {
		return nil, response.Unauthorized(err.Error())
	}
	return nil, nil
}

func decodeCreateUser(ctx *gin.Context) (any, error) {

	//params := ctx.Value("params").(map[string]string)

	if err := tokenVerify(ctx.Request.Header.Get("Authorization")); err != nil {
		return nil, response.Unauthorized(err.Error())
	}

	var req user.CreateReq
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil { //decodifica el body
		return nil, response.BadRequest(fmt.Sprintf("error decoding request: %v", err.Error()))
	}
	return req, nil
}

func encodeResponse(ctx *gin.Context, resp any) {

	//agrego el manejo del pkg response

	r := resp.(response.Response) //casteo el dato a un tipo Response
	ctx.Header("Content-Type", "application/json; charset=utf-8")
	ctx.JSON(r.StatusCode(), resp) //encodeo el response

}

func encodeError(ctx *gin.Context, err error) {

	ctx.Header("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response)   //casteo el dato a un tipo Response
	ctx.JSON(resp.StatusCode(), resp) //encodeo el response

}

func tokenVerify(token string) error {
	if os.Getenv("TOKEN") != token {
		return errors.New("invalid token")
	}
	return nil
}
