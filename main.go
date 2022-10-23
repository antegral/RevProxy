package main

import (
	"flag"
	"fmt"
	"net/http"

	Log "antegral.net/revproxy/src/Log"
	ProxyModule "antegral.net/revproxy/src/Proxy"
)

var (
	ListenPort string
	ProxyTo    string
	Password   string
)

func main() {
	Log.Init()

	flag.StringVar(&ListenPort, "listen", "", "Port to listen")
	flag.StringVar(&ProxyTo, "address", "", "Address to proxy")
	flag.StringVar(&Password, "password", "", "Header password (X-RevProxy-Token)")
	flag.Parse()

	Log.Info.Println("Starting RevProxy...")

	if len(ListenPort) <= 0 || len(ProxyTo) <= 0 {
		Log.Error.Panicln("Invaild listening port or address to proxy.")
	}

	if len(Password) <= 0 {
		Log.Info.Println("Password was not entered. Header authentication mode is not operational.")
	}

	Log.Info.Print("ListenPort: ", ListenPort, " / ", "ProxyTo: ", ProxyTo)

	if Proxy, err := ProxyModule.NewProxy(ProxyTo); err != nil {
	} else {
		http.HandleFunc("/", ProxyModule.RequestHandler(Proxy, Password))
		if err := http.ListenAndServe(fmt.Sprint(":", ListenPort), nil); err != nil {
			Log.Error.Panicln(err)
		}
	}
}
