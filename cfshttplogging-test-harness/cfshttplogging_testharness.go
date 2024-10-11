package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/connorsadler/go-utils/cfshttplogging"
)

//
// This test harness imports cfshttplogging so it can use: cfshttplogging.InstallLoggingRoundTripper
//

func main() {
	fmt.Println("Start")

	cfshttplogging.InstallLoggingRoundTripper(http.DefaultClient)

	// Example curl command:
	//   curl -v https://api.sampleapis.com/beers/ale/1 | jq
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
