package main

import (
	"log"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/pavozayac/scheduling/src/constraint-service/internal/application/graph"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		playground.Handler("GraphQL playground", "/query")(c.Writer, c.Request)
	})

	r.POST("/query", func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(r.Run())
}
