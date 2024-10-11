package cfshttplogging

import (
	"io"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: Old dummy test - remove this
func TestSampleFuncFromConnor(t *testing.T) {

	actualResult := SampleFuncFromConnor()

	assert.Regexp(t, "This is from cfshttplogging.go at.*", actualResult)
}

// Install the logging round tripper
// Call an external API, and check the result
// TODO: Assert the logging calls are shown
func TestLoggingRoundTripper(t *testing.T) {

	c := http.DefaultClient
	InstallLoggingRoundTripper(c)

	resp, err := http.Get("https://api.sampleapis.com/beers/ale/1")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Response body: %v", string(b))

	assert.True(t, strings.Contains(string(b), "Founders All Day IPA"))
}
