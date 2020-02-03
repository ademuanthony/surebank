package mid

import (
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"strings"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"merryworld/surebank/internal/platform/web"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

// Write uses the Writer part of gzipResponseWriter to write the output.
func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Compress() web.Middleware {
	f := func(after web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
			span, ctx := tracer.StartSpanFromContext(ctx, "internal.mid.compressor")
			defer span.Finish()
			// Check if the client can accept the gzip encoding.
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				return after(ctx, w, r, params)
			}
			// Set the HTTP header indicating encoding.
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(w)
			defer gz.Close()
			return after(ctx, gzipResponseWriter{Writer: gz, ResponseWriter: w}, r, params)
		}

		return h
	}
	return f
}
