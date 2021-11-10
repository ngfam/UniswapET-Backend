package main

import (
	"log"
	"net/http"
	"os"

	"internal/links"
	database "internal/pkg/db/mysql"
	"internal/users"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/docker/docker/daemon/links"
	"github.com/go-chi/chi"
	"github.com/ngfam/uniswap69420/auth"
	"github.com/ngfam/uniswap69420/graph"
	"github.com/ngfam/uniswap69420/graph/generated"
	"github.com/rs/cors"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	database.InitDB()
	database.Migrate()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	router := chi.NewRouter()
	router.Use(auth.Middleware())
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)
	handler := cors.AllowAll().Handler(router)
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func foo() {
	var x links.Link
	var y users.User
	log.Print(x.ID)
	log.Print(y.ID)
}
