package web

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type middleware func(http.Handler) http.HandlerFunc

func middlewaresChain(mid ...middleware) middleware {
	return func(nxt http.Handler) http.HandlerFunc {
		for i := len(mid) - 1; i >= 0; i-- {
			nxt = mid[i](nxt)
		}

		return nxt.ServeHTTP
	}
}

type logResponseWriter struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func newLogResponseWriter(rew http.ResponseWriter) *logResponseWriter {
	return &logResponseWriter{
		ResponseWriter: rew,
		status:         http.StatusOK,
		body:           bytes.NewBuffer(nil),
	}
}

func (w *logResponseWriter) Write(bdy []byte) (int, error) {
	w.body.Write(bdy)

	return w.ResponseWriter.Write(bdy)
}

func logRequestAndResponseMiddleware(nxt http.Handler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		uid := uuid.New().String()
		body, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(body))
		slog.Info("HTTP Request", "uid", uid, "method", req.Method, "path", req.RequestURI, "body", body)
		responseWriter := newLogResponseWriter(res)
		nxt.ServeHTTP(responseWriter, req)
		slog.Info("HTTP Response", "uid", uid, "status", responseWriter.status, "body", responseWriter.body.String())
	}
}
