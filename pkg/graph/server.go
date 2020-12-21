package graph

import (
	"github.com/99designs/gqlgen/graphql/handler"
)

func NewServer(r *Resolver, isProduction bool) *handler.Server {
	return handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: r}))
}
