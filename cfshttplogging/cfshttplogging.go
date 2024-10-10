package cfshttplogging

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// This type implements the http.RoundTripper interface
type LoggingRoundTripper struct {
	Proxied http.RoundTripper
}

const LoggingRoundTripper_Version = "v0.1"

func (lrt LoggingRoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	// Do "before sending requests" actions here.
	fmt.Printf("%v LoggingRoundTripper] >>> Sending request to %v\n", LoggingRoundTripper_Version, req.URL)
	// TODO: more info
	fmt.Printf("%v LoggingRoundTripper] >>> Headers: TODO\n", LoggingRoundTripper_Version)
	fmt.Printf("%v LoggingRoundTripper] >>> Body: %v\n", LoggingRoundTripper_Version, getRequestBody(req))

	// Send the request, get the response (or the error)
	res, e = lrt.Proxied.RoundTrip(req)

	// Handle the result.
	if e != nil {
		fmt.Printf("%v LoggingRoundTripper] <<< Error: %v\n", LoggingRoundTripper_Version, e)
	} else {
		fmt.Printf("%v LoggingRoundTripper] <<< Received %v response\n", LoggingRoundTripper_Version, res.Status)
	}
	// TODO: more info
	fmt.Printf("%v LoggingRoundTripper] <<< Body: %v\n", LoggingRoundTripper_Version, getResponseBody(res))

	return
}

func getRequestBody(req *http.Request) string {
	if req.Body == nil {
		return "[no body]"
	}
	buf := new(strings.Builder)
	_, err := io.Copy(buf, req.Body)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	req.Body.Close()

	// setup new request body for later consumption
	// TODO: Untested, but same as code below in getResponseBody
	req.Body = io.NopCloser(bytes.NewBuffer([]byte(buf.String())))

	return buf.String()
}

func getResponseBody(res *http.Response) string {
	if res.Body == nil {
		return "[no body]"
	}
	buf := new(strings.Builder)
	_, err := io.Copy(buf, res.Body)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	res.Body.Close()

	// setup new response body for later consumption
	res.Body = io.NopCloser(bytes.NewBuffer([]byte(buf.String())))

	return buf.String()
}

func InstallLoggingRoundTripper(c *http.Client) {
	origTransport := c.Transport
	if origTransport == nil {
		origTransport = http.DefaultTransport
	}
	c.Transport = LoggingRoundTripper{origTransport}
}

// func main() {
// 	fmt.Println("Start")

// 	// httpClient := &http.Client{
// 	// 	// Transport: LoggingRoundTripper{http.DefaultTransport},
// 	// }
// 	// httpClient.Get("https://example.com/")

// 	// See link: https://stackoverflow.com/questions/39527847/is-there-middleware-for-go-http-client
// 	// origTransport := http.DefaultClient.Transport
// 	// if origTransport == nil {
// 	// 	origTransport = http.DefaultTransport
// 	// }
// 	// http.DefaultClient.Transport = LoggingRoundTripper{origTransport}
// 	InstallLoggingRoundTripper(http.DefaultClient)

// 	// curl -v https://api.sampleapis.com/beers/ale/1 | jq
// 	resp, err := http.Get("https://api.sampleapis.com/beers/ale/1")
// 	if err != nil {
// 		log.Fatalf("Error: %v", err)
// 	}

// 	b, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	log.Printf("Response body: %v", string(b))

// 	fmt.Println("Done")
// }
