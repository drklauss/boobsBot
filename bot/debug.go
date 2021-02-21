package bot

import "sync"

var d debug

type debug struct {
	mx    sync.RWMutex
	value bool
}

// SetDebug changes debug state for bot.
func SetDebug(debug bool) {
	d.mx.Lock()
	d.value = debug
	d.mx.Unlock()
}

// Debug returns whether bot is in debug or not.
func Debug() bool {
	return d.value
}
