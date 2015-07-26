package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"net/http"
)

func RequestIdentification(serverId string) http.HandlerFunc {
	const (
		reqid = "x-request-id"
		dot   = "."
	)
	return func(w http.ResponseWriter, r *http.Request) {
		u1 := uuid.NewV4()
		var buffer bytes.Buffer
		buffer.WriteString(serverId)
		buffer.WriteString(dot)
		buffer.WriteString(u1.String())

		headers := w.Header()
		headers.Set(reqid, buffer.String())
	}
}

func GinRequestIdentification(server string) gin.HandlerFunc {
	req := RequestIdentification(server)
	return func(c *gin.Context) {
		req.ServeHTTP(c.Writer, c.Request)
	}
}
