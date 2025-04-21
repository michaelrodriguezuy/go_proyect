package transport

import (
	"context"
	"net/http"
)

type Transport interface {
	Server(
		endpoint Endpoint, //este va a ser nuestro controlador

		//estos serian los middlewares, que van a ser funciones que van a recibir un contexto y un request, y van a devolver un response y un error
		decode func(ctx context.Context, r *http.Request) (any, error), //este va a ser quien tome el request y lo decodifique, y se lo pase al controlador
		encode func(ctx context.Context, w http.ResponseWriter, response any) error, //este va a ser quien tome el response y lo codifique

		encodeError func(ctx context.Context, err error, w http.ResponseWriter) error, //este va a ser quien tome el error y lo codifique
	)
}

// este controlador va a recibir un contexto y un request (ya generado), y va a devolver un response y un error
type Endpoint func(ctx context.Context, request any) (response any, err error)

type transport struct {
	//este va a ser nuestro controlador
	w   http.ResponseWriter
	r   *http.Request
	ctx context.Context
}


//aqui, cuando se instancia la estructura tenemos el contexto, el request y el response
func NewTransport(w http.ResponseWriter, r *http.Request, ctx context.Context) Transport {
	return &transport{
		w:   w,
		r:   r,
		ctx: ctx,
	}
}

func (t *transport) Server(
	endpoint Endpoint,
	decode func(ctx context.Context, r *http.Request) (any, error),
	encode func(ctx context.Context, w http.ResponseWriter, response any) error,
	encodeError func(ctx context.Context, err error, w http.ResponseWriter) error,
) {
	
	//DECODIFICA EL REQUEST
	data, err := decode(t.ctx, t.r) //decodifica el request y lo guarda en la variable data
	if err != nil {	
		encodeError(t.ctx, err, t.w) //si hay un error, lo codifica y lo devuelve
		return
	}

	//PASA EL REQUEST DECODIFICADO AL CONTROLADOR, Y ESTE DEVUELVE UN RESPONSE
 	resp, err := endpoint(t.ctx, data) //llama al controlador y le pasa el contexto y el request decodificado
	if err != nil {
		encodeError(t.ctx, err, t.w) 
		return
	}

	//CODIFICA EL RESPONSE
	if err:= encode(t.ctx, t.w, resp); err != nil { //si no hay error, lo codifica y lo devuelve
		encodeError(t.ctx, err, t.w)
		return
	}

}