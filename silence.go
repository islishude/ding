package ding

import "sync"

var mux sync.RWMutex

var silence bool

// SetSilenceMode 设置是否静默模式，静默模式下，不发送任何消息
func SetSilenceMode(s bool) {
	mux.Lock()
	defer mux.Unlock()
	silence = s
}

func GetSilenceMode() bool {
	mux.RLock()
	defer mux.RUnlock()
	return silence
}
