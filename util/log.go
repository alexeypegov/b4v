package util

import (
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/urfave/negroni"
)

// LogMiddleware glog logger wrapper
type LogMiddleware struct {
	level glog.Level
}

// NewLogMiddleware creates new glog middleware
func NewLogMiddleware(level glog.Level) *LogMiddleware {
	return &LogMiddleware{
		level: level,
	}
}

func (m *LogMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	remoteAddr := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		remoteAddr = realIP
	}

	next(w, r)

	res := w.(negroni.ResponseWriter)
	if glog.V(m.level) {
		glog.Infof(`%s %s %s (remote=%s) (status=%d "%s") (size=%d) (duration=%s)`,
			r.Method,
			r.URL.Path,
			r.Proto,
			remoteAddr,
			res.Status(),
			http.StatusText(res.Status()),
			res.Size(),
			time.Since(start),
		)
	}
}
