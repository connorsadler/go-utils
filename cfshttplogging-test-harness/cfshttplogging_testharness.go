package cfshttplogging

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/connorsadler/go-utils/cfshttplogging"
)

//
// This file imports cfshttplogging so it can use: cfshttplogging.InstallLoggingRoundTripper
//

func main() {
	fmt.Println("Start")

	// httpClient := &http.Client{
	// 	// Transport: LoggingRoundTripper{http.DefaultTransport},
	// }
	// httpClient.Get("https://example.com/")

	// See link: https://stackoverflow.com/questions/39527847/is-there-middleware-for-go-http-client
	// origTransport := http.DefaultClient.Transport
	// if origTransport == nil {
	// 	origTransport = http.DefaultTransport
	// }
	// http.DefaultClient.Transport = LoggingRoundTripper{origTransport}
	// cfsutilspackage.InstallLoggingRoundTripper(http.DefaultClient)
	cfshttplogging.InstallLoggingRoundTripper(http.DefaultClient)

	// curl -v https://api.sampleapis.com/beers/ale/1 | jq
	resp, err := http.Get("https://api.sampleapis.com/beers/ale/1")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Response body: %v", string(b))

	fmt.Println("Done")
}
