package cfshttplogging

import (
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSampleFuncFromConnor(t *testing.T) {

	actualResult := SampleFuncFromConnor()

	assert.Regexp(t, "This is from cfshttplogging.go at.*", actualResult)
}

func TestInstallLoggingRoundTripper(t *testing.T) {

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

	assert.Equal(t, "TODO", "NOTDONEYET")
}
