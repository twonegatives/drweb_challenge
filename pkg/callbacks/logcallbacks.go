package callbacks

import (
	"fmt"
	"time"
)

type LogCallback struct {
	Content string
}

func (c *LogCallback) Invoke(args ...interface{}) {
	// TODO: change println to logging to file
	fmt.Println(fmt.Sprintf("[%s] %s", time.Now().Format(time.RFC3339), c.Content))
}
