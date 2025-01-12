# apperr

[![Go Report Card](https://goreportcard.com/badge/github.com/today2098/apperr)](https://goreportcard.com/report/github.com/today2098/apperr)
[![Go Reference](https://pkg.go.dev/badge/github.com/today2098/apperr.svg)](https://pkg.go.dev/github.com/today2098/apperr)

Custom error for WEB application by Golang.

[日本語](./README.ja.md)


## Features

- Holding HTTP status code and body
- Wrapping error
- Stacktrace
- Sentinel error
- Depend on only standard library


## Installation

```bash
go get -u github.com/today2098/apperr
```


## Usage

```go
import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/today2098/apperr"
)

type MyBody struct {
	Message string `json:"message"`
}

var _ (apperr.Body) = (*MyBody)(nil)

func (b *MyBody) Is(target apperr.Body) bool {
	return true
}

func (b MyBody) Clone() apperr.Body {
	return &b
}

func (b *MyBody) String() string {
	return b.Message
}

// Sentinel error.
var (
	ErrNotFound            = apperr.New(404, &MyBody{Message: "Not Found"})
	ErrInternalServerError = apperr.New(500, &MyBody{Message: "Internal Server Error"})
)

func Service() error {
	return ErrNotFound.Wrap("something error")
}

func Controller() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := Service(); err != nil {
			var appErr *apperr.Error
			if !errors.As(err, &appErr) {
				appErr = ErrInternalServerError.WrapPrefix(err, "controller")
			}

			log.Printf("%v\n", appErr)               // apperr(404): Not Found; something error
			log.Printf("%v\n", appErr.StackFrames()) // Print the stack frames.

			ctx.JSON(appErr.StatusCode, appErr.Body)
			return
		}

		ctx.Status(200)
	}
}

func main() {
	r := gin.Default()
	r.GET("/", Controller())
	r.Run()
}
```
