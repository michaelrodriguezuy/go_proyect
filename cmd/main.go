package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/michaelrodriguezuy/go_proyect/internal/user"
	"github.com/michaelrodriguezuy/go_proyect/pkg/bootstrap"
	"github.com/michaelrodriguezuy/go_proyect/pkg/handler"
)

func main() {
	server := http.NewServeMux()

	db := bootstrap.NewDB()
	logger := bootstrap.NewLogger()

	//simulo el pasaje por las distintas capas

	repo := user.NewRepository(db, logger)
	service := user.NewService(logger, repo)
	ctx := context.Background() //opcional, por si tenemos que pasar informacion a las diferentes capas

	handler.NewUserHTTPServer(ctx, server, user.NewEndpoint(ctx, service))

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", server))
}
