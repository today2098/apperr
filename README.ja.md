# apperr

[![Go Report Card](https://goreportcard.com/badge/github.com/today2098/apperr)](https://goreportcard.com/report/github.com/today2098/apperr)
[![Go Reference](https://pkg.go.dev/badge/github.com/today2098/apperr.svg)](https://pkg.go.dev/github.com/today2098/apperr)

Go による WEB アプリケーションの実装での利用を想定したカスタムエラーです．

[English](./README.md)


## 特徴

- HTTPステータスコードとボディ内容を保持
- エラーのラップ (wrap)
- スタックトレース
- センチネルエラー
- 標準ライブラリのみに依存


## インストール

```bash
go get -u github.com/today2098/apperr
```


## 使い方

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
