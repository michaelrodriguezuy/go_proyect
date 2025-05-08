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

	_=godotenv.Load() //carga las variables de entorno del archivo .env

	server := http.NewServeMux()

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

	handler.NewUserHTTPServer(ctx, server, user.NewEndpoint(ctx, service))

	port := os.Getenv("PORT")
	fmt.Println("Server is running on port ", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s",port), server))
}
