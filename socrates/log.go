package socrates

import (
	"fmt"
	"log"
	"os"
)


var logger = log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags)
var Verbose = false;

func logf(f string, v ...interface{}) {
	if Verbose {
		logger.Output(2, fmt.Sprintf(f, v...))
	}
}

type logHelper struct {
	prefix string
}

func (l *logHelper) Write(p []byte) (n int, err error) {
	if Verbose {
		logger.Printf("%s%s\n", l.prefix, p)
		return len(p), nil
	}
	return
}

func newLogHelper(prefix string) *logHelper {
	return &logHelper{prefix}
}
