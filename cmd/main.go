package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/michaelrodriguezuy/go_proyect/internal/user"
	"github.com/michaelrodriguezuy/go_proyect/pkg/bootstrap"
	"github.com/michaelrodriguezuy/go_proyect/pkg/handler"
)

func main() {

	_ = godotenv.Load() //carga las variables de entorno del archivo .env

	db, err := bootstrap.NewDBConnection()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	if condition := db.Ping(); condition != nil {
		log.Fatal(condition)
	}

	logger := bootstrap.NewLogger()
	//simulo el pasaje por las distintas capas
	repo := user.NewRepository(db, logger)
	service := user.NewService(logger, repo)

	ctx := context.Background() //opcional, por si tenemos que pasar informacion a las diferentes capas

	h := handler.NewUserHTTPServer(user.NewEndpoint(ctx, service))

	port := os.Getenv("PORT")
	fmt.Println("Server is running on port ", port)

	address := fmt.Sprintf("127.0.0.1:%s", port)

	srv := &http.Server{
		Handler: accessControl(h),
		Addr:    address,
	}

	log.Fatal(srv.ListenAndServe())
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")                                                                                                                            // Permitir cualquier origen
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")                                                                                      // Permitir métodos
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization, Cache-Control, DNT, If-Modified-Since, Keep-Alive, Origin, User-Agent,X-Requested-With") // Permitir encabezados

		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}
