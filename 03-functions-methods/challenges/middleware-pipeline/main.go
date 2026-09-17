package main

import (
	"fmt"
)

type Request struct {
	Path      string
	AuthToken string
	Body      string
}

type Response struct {
	StatusCode int
	Body       string
}

type HandlerFunc func(req Request) Response

type Middleware func(next HandlerFunc) HandlerFunc

// LoggingMiddleware logs incoming requests and response status codes.
func LoggingMiddleware(logger func(string)) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(req Request) Response {
			if logger != nil {
				logger(fmt.Sprintf("--> %s", req.Path))
			}
			res := next(req)
			if logger != nil {
				logger(fmt.Sprintf("<-- %s [%d]", req.Path, res.StatusCode))
			}
			return res
		}
	}
}

// AuthMiddleware verifies the auth token before delegating.
func AuthMiddleware(validToken string) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(req Request) Response {
			if req.AuthToken != validToken {
				return Response{
					StatusCode: 401,
					Body:       "Unauthorized",
				}
			}
			return next(req)
		}
	}
}

// RecoveryMiddleware intercepts downstream panics using defer and recover.
func RecoveryMiddleware(panicLogger func(any)) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(req Request) (res Response) {
			defer func() {
				if r := recover(); r != nil {
					if panicLogger != nil {
						panicLogger(r)
					}
					res = Response{
						StatusCode: 500,
						Body:       "Internal Server Error",
					}
				}
			}()
			return next(req)
		}
	}
}

// Chain wraps the final handler with all middlewares in execution order.
func Chain(middlewares ...Middleware) func(HandlerFunc) HandlerFunc {
	return func(final HandlerFunc) HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}

func main() {
	echoHandler := func(req Request) Response {
		return Response{StatusCode: 200, Body: "Echo: " + req.Body}
	}

	pipeline := Chain(
		LoggingMiddleware(func(msg string) { fmt.Println("[LOG]", msg) }),
		AuthMiddleware("secret123"),
	)(echoHandler)

	// Valid request
	res := pipeline(Request{Path: "/api/hello", AuthToken: "secret123", Body: "Hello Go"})
	fmt.Printf("Response: %d, %s\n", res.StatusCode, res.Body)

	// Unauthorized request
	res2 := pipeline(Request{Path: "/api/hello", AuthToken: "badtoken", Body: "Hello Go"})
	fmt.Printf("Response: %d, %s\n", res2.StatusCode, res2.Body)
}
