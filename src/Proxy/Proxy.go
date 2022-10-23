package Proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"antegral.net/revproxy/src/Log"
)

func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	if url, err := url.Parse(targetHost); err != nil {
		return nil, err
	} else {
		return httputil.NewSingleHostReverseProxy(url), nil
	}
}

func RequestHandler(proxy *httputil.ReverseProxy, Password string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		PassHeader := r.Header.Get("X-RevProxy-Token")
		if len(Password) <= 0 {
			proxy.ServeHTTP(w, r)
		}
		if PassHeader != Password {
			w.WriteHeader(403)
			w.Write([]byte("Forbidden Page"))
			Log.Warn.Print("[", r.Method, "/FAIL] ", r.RemoteAddr, " > ", r.RequestURI)
		} else {
			Log.Info.Print("[", r.Method, "/OK] ", r.RemoteAddr, " > ", r.RequestURI)
			proxy.ServeHTTP(w, r)
		}
	}
}
