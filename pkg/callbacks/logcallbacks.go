package callbacks

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

type LogCallback struct {
	Content string
}

func (c *LogCallback) Invoke(args ...interface{}) {
	log.Info(fmt.Sprintf("[%s] %s", time.Now().Format(time.RFC3339), c.Content))
}
