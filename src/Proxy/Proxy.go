package Proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"antegral.net/revproxy/src/Log"
)

func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	Log.Verbose.Print("NewProxy > targetHost: ", targetHost)
	if url, err := url.Parse(targetHost); err != nil {
		return nil, err
	} else {
		Log.Verbose.Print("NewProxy > url.path: ", url.Path)
		return httputil.NewSingleHostReverseProxy(url), nil
	}
}

func RequestHandler(proxy *httputil.ReverseProxy, Password string, ProxyUrl *url.URL) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		PassHeader := r.Header.Get("X-RevProxy-Token")

		r.Header.Set("Host", ProxyUrl.Hostname())
		r.Host = ProxyUrl.Hostname()

		if len(Password) <= 0 {
			Log.Info.Print("[", r.Method, "/OK] ", r.RemoteAddr, " > ", r.RequestURI)
			proxy.ServeHTTP(w, r)
		} else if PassHeader != Password {
			Log.Warn.Print("[", r.Method, "/FAIL] ", r.RemoteAddr, " > ", r.RequestURI)
			w.Write([]byte("Forbidden Page"))
		} else {
			Log.Info.Print("[", r.Method, "/OK] ", r.RemoteAddr, " > ", r.RequestURI)
			proxy.ServeHTTP(w, r)
		}
	}
}
