package context

import (
	"github.com/gin-gonic/gin"
	"io"
)

type Context struct {
	RemoteAddress string
	Host string
	Method string
	RequestUri string
	Header map[string][]string
	Params gin.Params
	Body io.ReadCloser
}

func NewContext(context *gin.Context) Context {
	return Context{
		RemoteAddress: context.Request.RemoteAddr,
		Host: context.Request.Host,
		Method: context.Request.Method,
		RequestUri: context.Request.RequestURI,
		Params: context.Params,
		Body: context.Request.Body,
	}
}
