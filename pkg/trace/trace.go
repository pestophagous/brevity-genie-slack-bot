package trace

import (
	"fmt"
	"log"
	"strings"
)

var enabledTraces map[string]int // The 'int' is currently ignored. Think of this as a set.

func init() {
	enabledTraces = make(map[string]int)
}

func EnableTraces(traces []string) {
	for _, t := range traces {
		enabledTraces[strings.ToLower(t)] = 1
	}
}

func Tracef(mask string, format string, v ...interface{}) {
	if _, ok := enabledTraces[strings.ToLower(mask)]; ok {
		log.Printf("[TRACE][%s] %s", mask, fmt.Sprintf(format, v...))
	}
}

func Trace(mask string, v ...interface{}) {
	if _, ok := enabledTraces[strings.ToLower(mask)]; ok {
		log.Print(fmt.Sprintf("[TRACE][%s] %s", mask, fmt.Sprint(v...)))
	}
}
