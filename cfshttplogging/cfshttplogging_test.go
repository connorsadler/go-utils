package cfshttplogging

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSampleFuncFromConnor(t *testing.T) {

	actualResult := SampleFuncFromConnor()

	assert.Regexp(t, "This is from cfshttplogging.go at.*", actualResult)
}
