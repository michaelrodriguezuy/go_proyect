package transport

import "github.com/gin-gonic/gin"

func GinServer(
	endpoint Endpoint,
	decode func(ctx *gin.Context) (any, error),
	encode func(ctx *gin.Context, response any),
	encodeError func(ctx *gin.Context, err error)) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		//DECODIFICA EL REQUEST
		data, err := decode(ctx) //decodifica el request y lo guarda en la variable data
		if err != nil {
			encodeError(ctx, err) //si hay un error, lo codifica y lo devuelve
			return
		}

		//PASA EL REQUEST DECODIFICADO AL CONTROLADOR, Y ESTE DEVUELVE UN RESPONSE
		resp, err := endpoint(ctx.Request.Context(), data) //llama al controlador y le pasa el contexto y el request decodificado
		if err != nil {
			encodeError(ctx, err)
			return
		}

		encode(ctx, resp)
	}
}
