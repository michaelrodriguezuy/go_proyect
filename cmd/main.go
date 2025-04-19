package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/michaelrodriguezuy/go_proyect/internal/domain"
	"github.com/michaelrodriguezuy/go_proyect/internal/user"
)

func main() {
	server := http.NewServeMux()

	db := user.DB{
		Users: []domain.User{
			{ID: 1, FirstName: "John", LastName: "Doe", Age: 30},
			{ID: 2, FirstName: "Jane", LastName: "Smith", Age: 25},
		},
		MaxUserID: 2,
	}

	//simulo el pasaje por las distintas capas
	logger := log.New(os.Stdout, "user: ", log.LstdFlags|log.Lshortfile)
	repo := user.NewRepository(db, logger)
	service := user.NewService(logger, repo)
	ctx := context.Background() //opcional, por si tenemos que pasar informacion a las diferentes capas

	server.HandleFunc("/users", user.NewEndpoint(ctx, service))

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", server))
}
