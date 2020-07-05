package gourd

import (
	"log"
	"net/http"
)

func Recovery() HandlerFunc {
	return func(context *Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("%s\n\n", err)
				context.WriteString(http.StatusInternalServerError, "500 Internal Server Error")
			}
		}()
		context.Next()
	}
}
