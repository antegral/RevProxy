package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"

	Log "antegral.net/revproxy/src/Log"
	ProxyModule "antegral.net/revproxy/src/Proxy"
)

var (
	ListenPort  string
	ProxyTo     string
	Password    string
	LoggingMode int
)

func main() {
	flag.StringVar(&ListenPort, "listen", "", "Port to listen")
	flag.StringVar(&ProxyTo, "address", "", "Address to proxy")
	flag.StringVar(&Password, "password", "", "Header password (X-RevProxy-Token)")
	flag.IntVar(&LoggingMode, "logging-mode", 3, "Logging mode")
	flag.Parse()

	Log.Init(LoggingMode)

	Log.Info.Println("Starting RevProxy...")

	if len(ListenPort) <= 0 || len(ProxyTo) <= 0 {
		Log.Error.Panicln("Invaild listening port or address to proxy.")
	}

	if len(Password) <= 0 {
		Log.Info.Println("Password was not entered. Header authentication mode is not operational.")
	}

	Log.Info.Print("ListenPort: ", ListenPort, " / ", "ProxyTo: ", ProxyTo)

	ProxyUrl, err := url.Parse(ProxyTo)
	if err != nil {
		Log.Error.Panicln(err)
	} else {
		Log.Verbose.Print("Host: ", ProxyUrl.Hostname())
	}

	if Proxy, err := ProxyModule.NewProxy(ProxyTo); err != nil {
		Log.Error.Panicln(err)
	} else {
		http.HandleFunc("/", ProxyModule.RequestHandler(Proxy, Password, ProxyUrl))
		if err := http.ListenAndServe(fmt.Sprint(":", ListenPort), nil); err != nil {
			Log.Error.Panicln(err)
		}
	}
}
