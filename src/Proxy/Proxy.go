package Proxy

import (
	"crypto/subtle"
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

func RequestHandler(proxy *httputil.ReverseProxy, password string, proxyUrl *url.URL) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		inputToken := r.Header.Get("X-RevProxy-Token")
		isCorrectPassword := password == "" || subtle.ConstantTimeCompare([]byte(inputToken), []byte(password)) == 1

		r.Header.Set("Host", proxyUrl.Hostname())
		r.Host = proxyUrl.Hostname()

		if !isCorrectPassword {
			Log.Warn.Print("[", r.Method, "/FAIL] ", r.RemoteAddr, " > ", r.RequestURI)
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Forbidden Page"))
		}

		Log.Info.Print("[", r.Method, "/OK] ", r.RemoteAddr, " > ", r.RequestURI)
		proxy.ServeHTTP(w, r)
	}
}
