/* ProxyGet
 */

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var proxyString string
var rawURL string
var debug = false

func main() {
	if len(os.Args) != 3 && len(os.Args) != 4 {
		fmt.Println("Usage: ", os.Args[0], "http://proxy-host:port http://host:port/page")
		os.Exit(1)
	}
	proxyString = os.Args[1]
	rawURL = os.Args[2]

	if len(os.Args) == 4 && os.Args[3] == DebugReqIdentifier {
		debug = true
	}

	fmt.Printf("%s %s\n", proxyString, rawURL)

	handler := NewReverseProxyViaProxy(rawURL, proxyString)

	http.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(":8080", nil))

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
