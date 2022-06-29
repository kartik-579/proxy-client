/* ProxyGet
 */

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

var proxyString string
var rawURL string

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: ", os.Args[0], "http://proxy-host:port http://host:port/page")
		os.Exit(1)
	}
	proxyString = os.Args[1]
	rawURL = os.Args[2]

	fmt.Printf("%s %s\n", proxyString, rawURL)

	handler := NewReverseProxyViaProxy(rawURL, proxyString)

	http.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func NewReverseProxyViaProxy(target string, proxy string) func(w http.ResponseWriter, r *http.Request) {
	targetURL, err := url.Parse(target)
	checkError(err)

	proxyURL, err := url.Parse(proxy)
	checkError(err)

	transport := &http.Transport{
		Proxy:             http.ProxyURL(proxyURL),
		DisableKeepAlives: false,
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(targetURL)
	reverseProxy.Transport = transport
	reverseProxy.Director = func(req *http.Request) {
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.Host = targetURL.Host
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		reverseProxy.ServeHTTP(w, r)
	}
}

func checkError(err error) {
	if err != nil {
		if err == io.EOF {
			return
		}
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
