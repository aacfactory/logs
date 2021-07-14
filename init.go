package logs

import (
	"fmt"
)

func init() {
	buildErr := build()
	if buildErr != nil {
		panic(fmt.Errorf("logs build failed, %v", buildErr))
	}
}
