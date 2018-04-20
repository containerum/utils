package swaggerutil

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
)

type DocFlavor int

const (
	FlavorRedoc DocFlavor = iota
	FlavorSwagger
)

func SwaggerUI(rd io.Reader, flavor DocFlavor) (gin.HandlerFunc, error) {
	data, err := ioutil.ReadAll(rd)
	if err != nil {
		return nil, err
	}

	spec, err := loads.Analyzed(json.RawMessage(data), "")
	if err != nil {
		return nil, err
	}

	handler := http.NotFoundHandler()
	switch flavor {
	case FlavorRedoc:
		handler = middleware.Redoc(middleware.RedocOpts{
			BasePath: "/",
		}, handler)
		fallthrough
	case FlavorSwagger:
		handler = middleware.Spec("/", spec.Raw(), handler)
	default:
		return nil, errors.New("invalid flavor type")
	}

	return func(ctx *gin.Context) {
		handler.ServeHTTP(ctx.Writer, ctx.Request)
	}, nil
}
