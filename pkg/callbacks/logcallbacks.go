package callbacks

import (
	"fmt"
	"log"
	"time"
)

type LogCallback struct {
	Content string
}

func (c *LogCallback) Invoke(args ...interface{}) {
	log.Println(fmt.Sprintf("[%s] %s", time.Now().Format(time.RFC3339), c.Content))
}
