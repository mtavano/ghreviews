package graph

import (
	"context"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gorilla/websocket"
	"github.com/mtavano/ghreviews"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func NewServer(r *Resolver, isProduction bool) *handler.Server {
	srv := handler.New(NewExecutableSchema(Config{Resolvers: r}))

	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 5 * time.Second,
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New(1000))

	if !isProduction {
		r.logger.Debugln("instropectioonnn")
		srv.Use(extension.Introspection{})
	}
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})
	srv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		if err, ok := e.(ghreviews.ServiceError); ok {
			extensions := map[string]interface{}{"code": err.Code}

			if !isProduction {
				extensions["error"] = err.Unwrap().Error()
			}

			return &gqlerror.Error{
				Message:    err.Error(),
				Path:       graphql.GetFieldContext(ctx).Path(),
				Extensions: extensions,
			}
		}

		return graphql.DefaultErrorPresenter(ctx, e)
	})

	return srv
}
