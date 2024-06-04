package raiden

import (
	"context"
	"encoding/json"

	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel/trace"
)

type (
	// The `Context` interface defines a set of methods that can be implemented by a struct to provide a
	// context for handling HTTP requests in the Raiden framework.
	Context interface {
		Ctx() context.Context
		SetCtx(ctx context.Context)

		Config() *Config

		SendRpc(Rpc) error
		ExecuteRpc(Rpc) (any, error)

		SendJson(data any) error
		SendError(message string) error
		SendErrorWithCode(statusCode int, err error) error

		RequestContext() *fasthttp.RequestCtx

		Span() trace.Span
		SetSpan(span trace.Span)

		Tracer() trace.Tracer

		Write(data []byte)
		WriteError(err error)
	}

	// The `Ctx` struct is a struct that implements the `Context` interface in the Raiden framework. It
	// embeds the `context.Context` and `*fasthttp.RequestCtx` types, which provide the context and request
	// information for handling HTTP requests. Additionally, it has fields for storing the configuration
	// (`config`), span (`span`), and tracer (`tracer`) for tracing and monitoring purposes.
	Ctx struct {
		context.Context
		*fasthttp.RequestCtx
		config *Config
		span   trace.Span
		tracer trace.Tracer
	}
)

func (c *Ctx) Config() *Config {
	return c.config
}

func (c *Ctx) SendRpc(rpc Rpc) error {
	rs, err := c.ExecuteRpc(rpc)
	if err != nil {
		return err
	}

	return c.SendJson(rs)
}

func (c *Ctx) ExecuteRpc(rpc Rpc) (any, error) {
	return ExecuteRpc(c, rpc)
}

func (c *Ctx) RequestContext() *fasthttp.RequestCtx {
	return c.RequestCtx
}

func (c *Ctx) Span() trace.Span {
	return c.span
}

func (c *Ctx) SetSpan(span trace.Span) {
	c.span = span
}

func (c *Ctx) Tracer() trace.Tracer {
	return c.tracer
}

func (c *Ctx) Ctx() context.Context {
	return c.Context
}

func (c *Ctx) SetCtx(ctx context.Context) {
	c.Context = ctx
}

// The `SendJson` function is a method of the `Ctx` struct in the Raiden framework. It is responsible
// for sending a JSON response to the client.
func (c *Ctx) SendJson(data any) error {
	c.Response.Header.SetContentType("application/json")
	byteData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	c.Write(byteData)
	return nil
}

func (c *Ctx) SendError(message string) error {
	return &ErrorResponse{
		Message:    message,
		StatusCode: fasthttp.StatusInternalServerError,
	}
}

func (c *Ctx) SendErrorWithCode(statusCode int, err error) error {
	return &ErrorResponse{
		Message:    err.Error(),
		StatusCode: statusCode,
		Code:       fasthttp.StatusMessage(statusCode),
	}
}

// The `WriteError` function is a method of the `Ctx` struct in the Raiden framework. It is responsible
// for writing an error response to the HTTP response body.
func (c *Ctx) WriteError(err error) {
	c.Response.Header.SetContentType("application/json")
	if errResponse, ok := err.(*ErrorResponse); ok {
		responseByte, errMarshall := json.Marshal(errResponse)
		if errMarshall == nil {
			c.Response.SetStatusCode(errResponse.StatusCode)
			c.Response.AppendBody(responseByte)
			return
		}
		err = errMarshall
	}
	c.Response.SetStatusCode(fasthttp.StatusInternalServerError)
	c.Response.AppendBodyString(err.Error())
}

// The `Write` function is a method of the `Ctx` struct in the Raiden framework. It is responsible for
// writing the response body to the HTTP response.
func (c *Ctx) Write(data []byte) {
	c.Response.SetStatusCode(fasthttp.StatusOK)
	c.Response.AppendBody(data)
}
