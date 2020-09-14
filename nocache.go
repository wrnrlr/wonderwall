package wonderwall

import (
	"net/http"
	"time"
)

var epoch = time.Unix(0, 0).Format(time.RFC1123)

var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

var etagHeaders = []string{
	"ETag",
	"If-Modified-Since",
	"If-Match",
	"If-None-Match",
	"If-Range",
	"If-Unmodified-Since",
}

func noCacheMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, v := range etagHeaders {
			if w.Header().Get(v) != "" {
				w.Header().Del(v)
			}
		}
		// Set our noCachMiddleware headers
		for k, v := range noCacheHeaders {
			w.Header().Set(k, v)
		}
		next.ServeHTTP(w, r)
	}
}
