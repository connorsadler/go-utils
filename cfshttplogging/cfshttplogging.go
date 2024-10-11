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

func logMessage(msg string) {
	fmt.Printf("[%v LoggingRoundTripper] %v\n", LoggingRoundTripper_Version, msg)
}

func (lrt LoggingRoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	// Do "before sending requests" actions here.
	logMessage(fmt.Sprintf(">>> Sending request to %v", req.URL))
	// TODO: more info
	logMessage(fmt.Sprintf(">>> Headers: TODO"))
	logMessage(fmt.Sprintf(">>> Body: %v", getRequestBody(req)))

	// Send the request, get the response (or the error)
	res, e = lrt.Proxied.RoundTrip(req)

	// Handle the result.
	if e != nil {
		logMessage(fmt.Sprintf("<<< Error: %v", e))
	} else {
		logMessage(fmt.Sprintf("<<< Received %v response", res.Status))
	}
	// TODO: more info
	logMessage(fmt.Sprintf("<<< Body: %v", getResponseBody(res)))

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
