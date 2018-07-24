package callbacks

import (
	"fmt"
	"time"
)

type LogCallback struct {
	Content string
}

func (c *LogCallback) Invoke(args ...interface{}) {
	fmt.Println(fmt.Sprintf("[%s] %s"), time.Now().Format(time.RFC3339), c.Content)
}
