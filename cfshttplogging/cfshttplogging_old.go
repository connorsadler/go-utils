package cfshttplogging

import (
	"fmt"
	"time"
)

func SampleFuncFromConnor() string {
	t := time.Now().UTC()
	formatted := fmt.Sprintf("%d%02d%02d_%02d%02d%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	return fmt.Sprintf("This is from cfshttplogging.go at %s", formatted)
}
