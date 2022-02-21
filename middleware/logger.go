package middleware

import (
	"log"
	"time"

	"github.com/izwzhang/chomolungma"
)

func Logger() chomolungma.HandlerFunc {
	return func(c *chomolungma.Context) {
		t := time.Now()
		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
